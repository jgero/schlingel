{
  description = "An application to share work";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, treefmt-nix }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
      # lib = pkgs.lib;
    in
    {
      formatter.${system} = treefmtEval.config.build.wrapper;
      packages.${system}.default = pkgs.buildGoModule {
        name = "schlingel";
        version = "0.0.1";
        # vendorHash = lib.fakeHash;
        vendorHash = "sha256-LGEKehYWz3QdpB87I1EL+xJ4NAcYjLaen+xvmNwCeDs=";
        preBuild = ''
          ${pkgs.templ}/bin/templ generate
        '';
        src = ./.;
      };
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          reflex
          templ
        ];
      };
    };

}
