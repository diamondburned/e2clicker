{ self, ... }:

{
  config,
  lib,
  pkgs,
  ...
}:

let
  e2clicker = config.services.e2clicker;
  backendFlags = "" + (lib.optionalString e2clicker.backend.debug "-vvv");
  frontendFlags = "";
in

with lib;
with builtins;

{
  options.services.e2clicker = {
    backend = {
      enable = mkEnableOption "e2clicker backend";

      host = mkOption {
        type = types.str;
        default = "127.0.0.1";
        description = "The host the backend should serve on.";
      };

      port = mkOption {
        type = types.int;
        description = "The port the backend should serve on.";
      };

      debug = mkOption {
        type = types.bool;
        default = false;
      };

      databaseURI = mkOption {
        type = types.str;
        description = "The URI of the database to use. Only PostgreSQL is supported.";
      };

      package = mkOption {
        type = types.package;
        default = self.packages.${pkgs.system}.e2clicker-backend;
      };
    };

    frontend = {
      enable = mkEnableOption "e2clicker frontend";

      host = mkOption {
        type = types.str;
        default = "127.0.0.1";
        description = "The host the frontend should serve on.";
      };

      port = mkOption {
        type = types.int;
        description = "The port the frontend should serve on.";
      };

      package = mkOption {
        type = types.package;
        default = self.packages.${pkgs.system}.e2clicker-frontend;
      };
    };
  };

  config = {
    systemd.services.e2clicker-backend = mkIf e2clicker.backend.enable {
      description = "e2clicker backend";
      after = [
        "network.target"
        "postgresql.service"
      ];
      wantedBy = [ "multi-user.target" ];
      environment = {
        HTTP_ADDRESS = "${e2clicker.backend.host}:${toString e2clicker.backend.port}";
        DATABASE_URI = "${e2clicker.backend.databaseURI}";
      };
      serviceConfig = {
        ExecStart = "${getExe e2clicker.backend.package} ${backendFlags}";
        Restart = "always";
        RestartSec = "5s";
        DynamicUser = true;
      };
    };

    systemd.services.e2clicker-frontend = mkIf e2clicker.frontend.enable {
      description = "e2clicker frontend";
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];
      environment = {
        HOST = "${e2clicker.frontend.host}";
        PORT = "${toString e2clicker.frontend.port}";
      };
      serviceConfig = {
        ExecStart = "${getExe e2clicker.frontend.package} ${frontendFlags}";
        Restart = "always";
        RestartSec = "5s";
        DynamicUser = true;
      };
    };
  };
}
