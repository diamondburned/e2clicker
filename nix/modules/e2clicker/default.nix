self:

{
  config,
  lib,
  pkgs,
  ...
}:

with lib;
with builtins;

let
  e2clicker = config.services.e2clicker;
  backendConfigFile = pkgs.writeText "e2clicker-backend.json" (builtins.toJSON e2clicker.backend);

  submoduleOptions = options: types.submodule { inherit options; };
  typeNullableSubmodule = { options }: types.nullOr (submoduleOptions options);
  typeJSONFile = { options }: types.either types.path (submoduleOptions options);
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
            description = "The HTTP timeout when making requests to notification servers.";
            type = types.str;
            default = "2m";
          };

          webPush = mkOption {
            description = ''
              The web push notification configuration. This contains the VAPID
              keys that are used to encrypt the notifications. Use `just
              generate-vapid` to generate the keys.
            '';
            type = types.nullOr (typeJSONFile {
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
            });
          };

          email = mkOption {
            description = ''
              The path to the file containing the email configuration in
              JSON. See `secrets/email-config.example.json` for an
              example.
            '';
            type = types.nullOr (typeJSONFile {
              options = {
                from = mkOption {
                  type = types.str;
                  description = "The email address to send notifications from.";
                };

                smtp = mkOption {
                  description = "The SMTP server configuration.";
                  type = types.submodule {
                    options = {
                      host = mkOption {
                        type = types.str;
                        description = "The SMTP server host.";
                      };
                      port = mkOption {
                        type = types.int;
                        description = "The SMTP server port.";
                      };
                      secure = mkOption {
                        type = types.bool;
                        default = true;
                        description = "Whether to use a secure connection.";
                      };
                      auth = mkOption {
                        type = types.submodule {
                          options = {
                            username = mkOption {
                              type = types.str;
                              description = "The SMTP server username.";
                            };
                            password = mkOption {
                              type = types.str;
                              description = "The SMTP server password.";
                            };
                          };
                        };
                      };
                    };
                  };
                };
              };
            });
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
        type = types.nullOr types.str;
        default = null;
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
      if e2clicker.frontend.socket then "/run/e2clicker-frontend.sock" else null;

    systemd.services.e2clicker-frontend = mkIf e2clicker.frontend.enable {
      description = "e2clicker frontend";
      environment = {
        NODE_ENV = "production";
        HOST_HEADER = "x-forwarded-host";
        PROTOCOL_HEADER = "x-forwarded-proto";
        # the frontend doesn't even handle POST requests.
        BODY_SIZE_LIMIT = "4K";
        IDLE_TIMEOUT = toString (1 * 60 * 60); # 1 hour
      };
      wantedBy = [ "multi-user.target" ];
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
      socketConfig = {
        ListenStream =
          if e2clicker.frontend.socket then
            "${e2clicker.frontend.socketPath}"
          else if e2clicker.frontend.host != null then
            "${e2clicker.frontend.host}:${toString e2clicker.frontend.port}"
          else
            "${toString e2clicker.frontend.port}";
        Accept = false;
        NoDelay = true;
        SocketMode = "0666";
        DirectoryMode = "0777";
      };
    };
  };
}
