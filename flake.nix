{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?rev=d04953086551086b44b6f3c6b7eeb26294f207da";
    flake-utils.url = "github:numtide/flake-utils";

    oapi-codegen = {
      url = "github:diamondburned/oapi-codegen?ref=migrate-to-libopenapi-ordered-fix";
      flake = false;
    };

    yaml-language-server = {
      url = "github:okybr/yaml-language-server";
      flake = false;
    };

    npmlock2nix = {
      url = "github:nix-community/npmlock2nix";
      flake = false;
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }@inputs:

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

        # Downgrade yaml-language-server to 1.15.0 to fix an issue with OpenAPI's v3.0.0 schema.
        yaml-language-server = pkgs.yaml-language-server.overrideAttrs (old: rec {
          version = "1.15.0-ajv-draft-04";
          src = inputs.yaml-language-server;
          doCheck = false;
          offlineCache = pkgs.fetchYarnDeps {
            yarnLock = "${src}/yarn.lock";
            hash = "sha256-thJ3aU52yCusfjBCD2QvLynwiM32lq0IT9WaNJjfu6E=";
          };
        });
      in
      {
        devShells.default = pkgs.mkShell {
          packages =
            with pkgs;
            with self.packages.${system};
            [
              just
              extra-container

              go
              gopls
              gotools # contains goimports
              go-tools # contains staticcheck
              moq

              (stub "npm")
              (stub "npx")
              nodejs
              nodePackages.pnpm

              sqlc
              pgformatter

              oapi-codegen
              yaml-language-server
              yq-go

              self.formatter.${system}
            ];

          shellHook = ''
            export PATH="$PATH:$(git rev-parse --show-toplevel)/node_modules/.bin"

            # Set up autocompletion for just in bash.
            [[ $SHELL == bash ]] && complete -W '$(just --summary)' just
          '';
        };

        packages = {
          oapi-codegen = pkgs.buildGoModule rec {
            pname = "oapi-codegen";
            src = inputs.oapi-codegen;
            version = builtins.substring 0 9 src.rev;
            subPackages = [ "cmd/oapi-codegen" ];
            doCheck = false;

            vendorHash = "sha256-aqjk+iAsO6rFqoqXJTMxeD2ZFR4FDrg+VbxGIE7XQzw=";
          };
        };

        formatter = pkgs.nixfmt-rfc-style;
      }
    );
}
