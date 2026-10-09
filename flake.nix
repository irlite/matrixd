{
  description = "matrixd development environment";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  outputs = {nixpkgs, ...}: let
    systems = [
      "x86_64-linux"
      "aarch64-linux"
    ];
    forAll = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
  in {
    devShells = forAll (pkgs: {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
          gotools
          delve
          gcc
          pkg-config
          golangci-lint
          golines
          gofumpt
          socat
        ];
      };
    });
  };
}
