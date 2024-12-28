{
  projectRootFile = "flake.nix";
  programs.nixpkgs-fmt.enable = true;
  programs = {
    gofmt.enable = true;
    templ.enable = true;
    typos.enable = true;
    prettier = {
      enable = true;
      includes = [ "*.md" "*.yaml" "*.yml" ];
    };
  };
}
