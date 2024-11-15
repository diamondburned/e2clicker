{
  config,
  lib,
  pkgs,
  self,
  ...
}:

with lib;
with builtins;

let
  e2clicker = config.services.e2clicker;
  backendConfigFile = pkgs.writeText "e2clicker-backend.json" (builtins.toJSON e2clicker.backend);

  submoduleOptions = options: types.submodule { inherit options; };
in

{
  options.services.e2clicker = {
    backend = {
      enable = mkEnableOption "e2clicker backend";

      api = mkOption {
        description = "configuration for the API server";
        type = submoduleOptions {
          listenAddress = mkOption {
            type = types.str;
            default = ":8080";
            description = "The address the API server should listen on.";
          };
        };
      };

      postgresql = mkOption {
        description = "configuration for the PostgreSQL database";
        type = submoduleOptions {
          databaseURI = mkOption {
            type = types.str;
            description = "The URI of the database to use.";
          };
        };
      };

      notification = mkOption {
        description = "configuration for the notification service";
        type = submoduleOptions {
          clientTimeout = mkOption {
            type = types.str;
            default = "2m";
            description = "The HTTP timeout when making requests to notification servers.";
          };

          webPushKeys = mkOption {
            type =
              let
                schema = types.submodule {
                  options = {
                    privateKey = mkOption {
                      type = types.str;
                      description = "The VAPID private key.";
                    };
                    publicKey = mkOption {
                      type = types.str;
                      description = "The VAPID public key.";
                    };
                  };
                };
              in
              types.nullOr (types.either types.path schema);
            description = "The path to the file containing the VAPID keys encoded in JSON";
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
        default = 8080;
        description = "The port the frontend should serve on.";
      };

      socket = mkOption {
        type = types.bool;
        default = false;
        description = ''
          Use a Unix socket instead of a TCP port.
          The path can be obtained from socketPath.
        '';
      };

      socketPath = mkOption {
        type = types.nullOr types.str;
        readOnly = true;
        description = ''
          The path to the socket the frontend will serve on.
          This is only not null if socket=true.
        '';
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

    services.e2clicker.frontend.socketPath =
      if e2clicker.frontend.socket then "/run/e2clicker-frontend/http.sock" else null;

    systemd.services.e2clicker-frontend = mkIf e2clicker.frontend.enable {
      description = "e2clicker frontend";
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];
      environment =
        (
          if e2clicker.frontend.socket then
            {
              SOCKET_PATH = "${e2clicker.frontend.socketPath}";
            }
          else
            {
              HOST = "${e2clicker.frontend.host}";
              PORT = "${toString e2clicker.frontend.port}";
            }
        )
        // {
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
        RuntimeDirectory = "e2clicker-frontend";
      };
    };

    systemd.sockets.e2clicker-frontend = mkIf e2clicker.frontend.enable {
      description = "e2clicker frontend socket";
      after = [ "network.target" ];
      wantedBy = [ "sockets.target" ];
      listenStreams =
        if e2clicker.frontend.socket then
          [ "${e2clicker.frontend.socketPath}" ]
        else
          [ "${e2clicker.frontend.host}:${toString e2clicker.frontend.port}" ];
    };
  };
}
