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
      url = "github:nix-community/gomod2nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "flake-utils";
      };
    };

    nixmod2go = {
      url = "github:diamondburned/nixmod2go";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    oapi-codegen = {
      url = "github:oapi-codegen/oapi-codegen/v2.6.0";
      flake = false;
    };

    yaml-language-server = {
      url = "github:redhat-developer/yaml-language-server";
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
              go = super.go_1_26;
              pnpm = super.pnpm_9;
            })

            (self: super: {
              oapi-codegen = pkgs.buildGoModule rec {
                pname = "oapi-codegen";
                src = inputs.oapi-codegen;
                version = builtins.substring 0 9 src.rev;
                subPackages = [ "cmd/oapi-codegen" ];
                doCheck = false;

                vendorHash = "sha256-vgSMGi0mnGX/Hwxu/XalIXLCbm/L4CwQfIf7DEJVk1E=";
              };

              pgformatter = import ./nix/pgformatter.nix { pkgs = super; };
            })
          ];
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

            nodejs
            pnpm

            sqlc
            pgformatter

            oapi-codegen
            redocly
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

            alias npm="echo Use pnpm instead of npm."
            alias npx="echo Use pnpx instead of npx."
          '';

          NO_UPDATE_NOTIFIER = "1";
        };

        packages = import ./nix/packages.nix { inherit pkgs self inputs; };

        formatter = pkgs.nixfmt;
      }
    ))
    // {
      nixosModules = {
        e2clicker = import ./nix/modules/e2clicker self;
        e2clicker-postgresql = import ./nix/modules/e2clicker-postgresql;
      };

      nixosConfigurations = builtins.listToAttrs (
        nixpkgs.lib.flip map
          [
            "x86_64-linux"
            "aarch64-linux"
          ]
          (system: {
            name = "dev-vm-${system}";
            value = nixpkgs.lib.nixosSystem {
              inherit system;
              specialArgs = inputs;
              modules = [
                inputs.nixos-shell.nixosModules.nixos-shell
                ./nix/dev/vm.nix
              ];
            };
          })
      );
    };
}
