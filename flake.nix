{
  description = "a demo of nix's capabilities";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        name = "nix-demo";

        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfreePredicate = pkg:
            builtins.elem (pkgs.lib.getName pkg) [ "terraform" ];
        };

        imageParams = { inherit pkgs name; };
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [ goose pgcli nodejs corepack go terraform ];
          GOOSE_DRIVER = "postgres";
          GOOSE_DBSTRING =
            "postgresql://postgres:password@localhost:8001/postgres?sslmode=disable";
          GOOSE_MIGRATION_DIR = ./db;
        };

        packages.db = import ./image-db.nix imageParams;
        packages.backend = import ./image-backend.nix imageParams;
        packages.frontend = import ./image-frontend.nix imageParams;
      });
}
