{ pkgs ? import <nixpkgs> {} }:

with pkgs;

mkShell {
  buildInputs = [go];

  shellHook = ''
    go get github.com/BattlesnakeOfficial/rules/cli/battlesnake
    export PATH=$HOME/go/bin:$PATH
  '';
}
