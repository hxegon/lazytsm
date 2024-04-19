{
  description = "A very basic flake";

  inputs = { nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable"; };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      devDeps = [ pkgs.just ];
    in {
      devShell.${system} =
        pkgs.mkShell { buildInputs = with pkgs; [ go ] ++ devDeps; };
    };
}
