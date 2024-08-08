{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        lib = pkgs.lib;

        pnpx =
          name:
          pkgs.writeShellScriptBin name ''
            exec pnpx -- ${name} "$@"
          '';

        stub =
          name:
          pkgs.writeShellScriptBin name ''
            echo "This command should not be run."
            exit 1
          '';
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            go-tools

            (stub "npm")
            (stub "npx")

            nodejs
            nodePackages.pnpm

            sqlc
            oapi-codegen
            yaml-language-server

            self.formatter.${system}
          ];

          shellHook = ''
            export PATH="$PATH:$(git rev-parse --show-toplevel)/node_modules/.bin"
          '';
        };

        formatter = pkgs.nixfmt-rfc-style;
      }
    );
}
