{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    nixos-shell = {
      url = "github:Mic92/nixos-shell";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };

    gomod2nix = {
      url = "github:obreitwi/gomod2nix?ref=fix/go_mod_vendor";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "flake-utils";
      };
    };

    nixmod2go = {
      url = "github:diamondburned/nixmod2go";
    };

    oapi-codegen = {
      url = "github:oapi-codegen/oapi-codegen/v2.4.1";
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
      nixmod2go,
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
              go = super.go_1_23;
              pnpm = super.pnpm_9;
            })

            (self: super: {
              oapi-codegen = pkgs.buildGoModule rec {
                pname = "oapi-codegen";
                src = inputs.oapi-codegen;
                version = builtins.substring 0 9 src.rev;
                subPackages = [ "cmd/oapi-codegen" ];
                doCheck = false;

                vendorHash = "sha256-bp5sFZNJFQonwfF1RjCnOMKZQkofHuqG0bXdG5Hf3jU=";
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

              pgformatter = import ./nix/pgformatter.nix { pkgs = super; };
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
            redocly-cli
            yaml-language-server
            yq-go

            nixmod2go.packages.${system}.default
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
        e2clicker = import ./nix/modules/e2clicker self;
        e2clicker-postgresql = import ./nix/modules/e2clicker-postgresql;
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
