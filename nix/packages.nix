{
  pkgs,
  self,
  inputs,
  ...
}@args:

let
  lib = pkgs.lib;
  hashes = with builtins; fromJSON (readFile ./hashes.json);

  src = self;
  version = if self ? "rev" then builtins.substring 0 9 self.rev else "devel";

  buildDeps = with pkgs; [
    just
    go
    moq
    nodejs
    pnpm
    sqlc
    oapi-codegen
    yq-go
  ];

  dist = pkgs.stdenv.mkDerivation rec {
    inherit src version;
    pname = "e2clicker-dist";

    nativeBuildInputs = buildDeps ++ [ pkgs.pnpm.configHook ];

    pnpmDeps = pkgs.pnpm.fetchDeps {
      inherit pname version src;
      hash = hashes.pnpmPackages;
    };

    goModules =
      (pkgs.buildGoApplication {
        inherit src version;
        pname = "e2clicker";
        modules = ./gomod2nix.toml;
      }).vendorEnv;

    # Needed for gomod2nix.
    # GO111MODULE = "on";
    # GOFLAGS = "-mod=vendor";

    buildPhase = ''
      runHook preBuild

      ln -s ${goModules} vendor

      just --no-deps build-backend
      just --no-deps build-frontend

      runHook postBuild
    '';

    installPhase = ''
      runHook preInstall

      cp -r dist $out

      # https://kit.svelte.dev/docs/adapter-node#deploying
      cp -r node_modules package.json $out/frontend/

      runHook postInstall
    '';
  };
in

rec {
  e2clicker-backend =
    pkgs.runCommandLocal "e2clicker-backend"
      {
        inherit version;
        meta = {
          description = "The e2clicker backend package";
          mainProgram = "e2clicker-backend";
        };
      }
      ''
        mkdir $out
        ln -s ${dist}/backend $out/bin
      '';

  e2clicker-frontend = pkgs.writeShellScriptBin "e2clicker-frontend" ''
    cd ${dist}/frontend
    exec ${lib.getExe pkgs.nodejs} .
  '';
}
