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

  frontend = pkgs.stdenv.mkDerivation (final: {
    inherit src version;
    pname = "e2clicker-frontend";

    nativeBuildInputs = buildDeps ++ [
      # This uses pnpmDeps.
      pkgs.pnpm.configHook
    ];

    buildInputs = with pkgs; [
      nodejs
    ];

    pnpmDeps = pkgs.pnpm.fetchDeps {
      inherit (final) pname version src;
      hash = hashes.pnpmPackages;
    };

    buildPhase = ''
      runHook preBuild

      just --no-deps build-frontend

      runHook postBuild
    '';

    installPhase = ''
      runHook preInstall

      mkdir $out
      cp -r dist $out/share
      # https://kit.svelte.dev/docs/adapter-node#deploying
      cp -r node_modules package.json $out/share/frontend/

      mkdir $out/bin
      {
        echo '#!/bin/bash' # will be fixed by stdenv
        echo "cd $out/share/frontend"
        echo 'exec '${lib.getExe pkgs.nodejs}' .'
      } > $out/bin/e2clicker-frontend
      chmod +x $out/bin/e2clicker-frontend

      runHook postInstall
    '';

    passthru = {
      assets = "${self}/share/frontend/client";
    };

    meta = {
      description = "The e2clicker frontend package";
      mainProgram = "e2clicker-frontend";
    };
  });

  backend = pkgs.buildGoApplication {
    inherit src version;
    pname = "e2clicker-backend";
    modules = ./gomod2nix.toml;

    nativeBuildInputs = buildDeps;

    buildPhase = ''
      runHook preBuild

      just --no-deps build-backend

      runHook postBuild
    '';

    installPhase = ''
      runHook preInstall

      mkdir $out
      cp -r dist $out/share

      mkdir $out/bin
      ln -s $out/share/backend/* $out/bin/

      runHook postInstall
    '';

    meta = {
      description = "The e2clicker backend package";
      mainProgram = "e2clicker-backend";
    };
  };
in

rec {
  e2clicker-frontend = frontend;
  e2clicker-backend = backend;
}
