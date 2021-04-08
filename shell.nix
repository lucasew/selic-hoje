{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gopls
    nodePackages.vercel
  ];
}
