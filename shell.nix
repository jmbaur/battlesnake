{ pkgs ? import <nixpkgs> {} }:

with pkgs;

mkShell {
  buildInputs = [go];

  shellHook = ''
    export PATH=$HOME/go/bin:$PATH

    if ! command -v battlesnake &>/dev/null
    then
      go get github.com/BattlesnakeOfficial/rules/cli/battlesnake
    fi
  '';
}
