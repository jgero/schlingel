{
  projectRootFile = "flake.nix";
  programs.nixpkgs-fmt.enable = true;
  programs.rustfmt.enable = true;
  programs.toml-sort.enable = true;
  programs.typos.enable = true;
  programs.yamlfmt.enable = true;
  programs.prettier.enable = true;
}
