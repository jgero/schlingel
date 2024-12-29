{
  description = "An application to share work";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-24.11";
    rust-overlay = { url = "github:oxalica/rust-overlay"; };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { nixpkgs, rust-overlay, treefmt-nix, ... }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ rust-overlay.overlays.default ];
      };
      treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
      rustWasm = (pkgs.rust-bin.stable.latest.default.override {
        extensions = [ "rust-src" ];
        targets = [ "wasm32-unknown-unknown" ];
      });
      rustPlatform = pkgs.makeRustPlatform {
        rustc = rustWasm;
        cargo = rustWasm;
      };
    in
    {
      formatter.${system} = treefmtEval.config.build.wrapper;
      checks.${system}.rust = pkgs.writeShellApplication {
        name = "cargo clippy";
        runtimeInputs = [ rustWasm pkgs.clippy ];
        text = ''
          cargo clippy
        '';
      };
      packages.${system}.default = rustPlatform.buildRustPackage rec {
        pname = "schlingel";
        version = "0.1.0";
        src = ./.;
        cargoLock = {
          lockFile = ./Cargo.lock;
        };
        nativeBuildInputs = with pkgs; [
          rustWasm
          binaryen
          wasm-pack
          cargo-leptos
        ];
        buildPhase = ''
          HOME=$TMPDIR cargo leptos build --release -vv
        '';
        installPhase = ''
          mkdir -p $out/bin
          ${pkgs.tree}/bin/tree .
          mv target/release/${pname} $out/bin/${pname}
          mv target/site/ $out/
        '';
      };
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          rustWasm
          cargo
          cargo-watch
          nodejs
          wasm-pack
          cargo-leptos
          cargo-generate
          binaryen
          clippy
        ];
      };
    };
}
