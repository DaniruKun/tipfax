{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
  ];

  shellHook = ''
    export SE_JWT_TOKEN=""
  '';
}
