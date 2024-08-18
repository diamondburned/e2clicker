{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?rev=d04953086551086b44b6f3c6b7eeb26294f207da";
    flake-utils.url = "github:numtide/flake-utils";

    nixos-shell = {
      url = "github:Mic92/nixos-shell";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "flake-utils";
      };
    };

    oapi-codegen = {
      url = "github:diamondburned/oapi-codegen?ref=migrate-to-libopenapi-ordered-fix";
      flake = false;
    };

    yaml-language-server = {
      url = "github:okybr/yaml-language-server";
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

    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            inputs.gomod2nix.overlays.default

            (self: super: {
              go = super.go_1_22;
              pnpm = super.pnpm_9;
            })

            (self: super: {
              oapi-codegen = pkgs.buildGoModule rec {
                pname = "oapi-codegen";
                src = inputs.oapi-codegen;
                version = builtins.substring 0 9 src.rev;
                subPackages = [ "cmd/oapi-codegen" ];
                doCheck = false;

                vendorHash = "sha256-aqjk+iAsO6rFqoqXJTMxeD2ZFR4FDrg+VbxGIE7XQzw=";
              };

              # Downgrade yaml-language-server to 1.15.0 to fix an issue with OpenAPI's v3.0.0 schema.
              yaml-language-server = super.yaml-language-server.overrideAttrs (old: rec {
                version = "1.15.0-ajv-draft-04";
                src = inputs.yaml-language-server;
                doCheck = false;
                offlineCache = super.fetchYarnDeps {
                  yarnLock = "${src}/yarn.lock";
                  hash = "sha256-thJ3aU52yCusfjBCD2QvLynwiM32lq0IT9WaNJjfu6E=";
                };
              });
            })
          ];
        };

        lib = pkgs.lib;

        stub =
          name:
          pkgs.writeShellScriptBin name ''
            echo "This command should not be run."
            exit 1
          '';

        hashes = {
          goModules = "sha256-5pWXRiZcnhk4N7wQYyiHGfOVhqUplyeiIJs9+PLn8fc=";
          pnpmPackages = "sha256-Qioex2l82EeJAQTaFQX/cT+ZFLJeT6ULI/+UfxaE9tk=";
        };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            just
            zellij

            go
            gopls
            gotools # contains goimports
            go-tools # contains staticcheck
            gomod2nix
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
            export PATH+=":$(git rev-parse --show-toplevel)/scripts"
            export PATH+=":$(git rev-parse --show-toplevel)/node_modules/.bin"

            # Set up autocompletion for just in bash.
            [[ $SHELL == bash ]] && complete -W '$(just --summary)' just
          '';

          NO_UPDATE_NOTIFIER = "1";
        };

        packages = import ./nix/packages.nix { inherit pkgs self inputs; };

        formatter = pkgs.nixfmt-rfc-style;
      }
    ))
    // {
      nixosModules = {
        e2clicker = import ./nix/modules/e2clicker.nix inputs;
        e2clicker-postgresql = import ./nix/modules/e2clicker-postgresql.nix inputs;
      };

      nixosConfigurations = {
        dev-vm = nixpkgs.lib.nixosSystem {
          system = "x86_64-linux";
          specialArgs = inputs;
          modules = [
            inputs.nixos-shell.nixosModules.nixos-shell
            ./nix/dev/vm.nix
          ];
        };
      };
    };
}
