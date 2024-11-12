{
  config,
  lib,
  pkgs,
  self,
  ...
}:

let
  e2clicker = config.services.e2clicker;
  backendConfigFile = pkgs.writeText "e2clicker-backend.json" (builtins.toJSON e2clicker.backend);
in

with lib;
with builtins;

{
  options.services.e2clicker = {
    backend = {
      enable = mkEnableOption "e2clicker backend";

      api = mkOption {
        description = "configuration for the API server";
        type = types.submodule {
          options = {
            listenAddress = mkOption {
              type = types.str;
              default = ":8080";
              description = "The address the API server should listen on.";
            };
          };
        };
      };

      postgresql = mkOption {
        description = "configuration for the PostgreSQL database";
        type = types.submodule {
          options = {
            databaseURI = mkOption {
              type = types.str;
              description = "The URI of the database to use.";
            };
          };
        };
      };

      debug = mkOption {
        description = "Enable debug logging and other debug features.";
        type = types.bool;
        default = false;
      };

      logFormat = mkOption {
        description = "The format of the log output.";
        type = types.enum [
          "color"
          "json"
          "text"
        ];
        default = "color";
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
      serviceConfig = {
        ExecStart = "${getExe e2clicker.backend.package} -c ${backendConfigFile}";
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
        HOST_HEADER = "x-forwarded-host";
        PROTOCOL_HEADER = "x-forwarded-proto";
        # the frontend doesn't even handle POST requests.
        BODY_SIZE_LIMIT = "4K";
        IDLE_TIMEOUT = toString (1 * 60 * 60); # 1 hour
      };
      serviceConfig = {
        ExecStart = "${getExe e2clicker.frontend.package}";
        Restart = "always";
        RestartSec = "5s";
        DynamicUser = true;
      };
    };

    systemd.sockets.e2clicker-frontend = mkIf e2clicker.frontend.enable {
      description = "e2clicker frontend socket";
      after = [ "network.target" ];
      wantedBy = [ "sockets.target" ];
      listenStreams = [ "${e2clicker.frontend.host}:${toString e2clicker.frontend.port}" ];
    };
  };
}
