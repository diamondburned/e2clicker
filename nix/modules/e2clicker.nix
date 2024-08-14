{
  config,
  lib,
  pkgs,
  self,
  ...
}:

let
  e2clicker = config.services.e2clicker;
in

with lib;
with builtins;

assert assertMsg (!e2clicker.frontend.enable) ("e2clicker frontend is not supported yet");

{
  options.services.e2clicker = {
    backend = {
      enable = mkEnableOption "e2clicker backend";

      httpAddress = mkOption {
        type = types.str;
        description = "The HTTP address the backend should serve on.";
      };

      databaseURI = mkOption {
        type = types.str;
        description = "The URI of the database to use. Only PostgreSQL is supported.";
      };
    };
    frontend = {
      enable = mkEnableOption "e2clicker frontend";
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
        HTTP_ADDRESS = e2clicker.backend.httpAddress;
        DATABASE_URI = e2clicker.backend.databaseURI;
      };
      serviceConfig = {
        ExecStart = getExe self.packages.${pkgs.system}.e2clicker-backend;
        Restart = "always";
        RestartSec = "5s";
        DynamicUser = true;
      };
    };
  };
}
