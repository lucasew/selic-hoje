{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gopls
    golangci-lint
    nodePackages.vercel
  ];
}
