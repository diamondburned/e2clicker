{
  config,
  lib,
  pkgs,
  ...
}@inputs:

let
  devFlags = lib.splitString " " (builtins.getEnv "E2CLICKER_DEVFLAGS");
  hasFlag = flag: lib.elem flag devFlags;
  noFlag = flag: !(hasFlag flag);

  root = ../..;
  rootPath = path: "${root}/${path}";
  rootPathOrNull =
    path: warn:
    if builtins.pathExists (rootPath path) then rootPath path else lib.warn "${warn}: at ${path}" null;
in

{
  imports = [
    inputs.self.nixosModules.e2clicker
    inputs.self.nixosModules.e2clicker-postgresql
  ];

  virtualisation = {
    diskSize = 4 * 1024;
    memorySize = 1024;

    forwardPorts = [
      {
        from = "host";
        host.port = 8000;
        guest.port = 80;
      }
    ];
  };

  systemd.extraConfig = ''
    DefaultStandardOutput=journal+console
  '';

  services.postgresql = {
    enable = true;
    enableJIT = true;
    settings = {
      log_connections = true;
      log_statement = "all";
    };
  };

  services.e2clicker = {
    frontend = {
      enable = noFlag "no-frontend";
      socket = true;
      trustProxy = true;
    };
    backend = {
      enable = noFlag "no-backend";
      debug = true;
      api = {
        listenAddress = ":36001";
      };
      postgresql = {
        databaseURI = "postgresql://e2clicker-backend@/e2clicker-backend";
      };
      notification = {
        email = rootPathOrNull "secrets/email-config.json" "No SMTP support because SMTP credentials missing";
        webPush = rootPathOrNull "secrets/vapid-keys.json" "No push notification support because VAPID keys missing";
        clientTimeout = "15s";
      };
    };
  };

  services.caddy = {
    enable = true;
    virtualHosts.":80".extraConfig =
      ""
      + (lib.optionalString (noFlag "no-frontend")) ''
        handle /_app/immutable* {
          header Cache-Control "public, immutable, max-age=31536000"
          file_server {
            root ${config.services.e2clicker.frontend.package.assets}
            precompressed gzip br
            pass_thru
          }
        }
        handle {
          reverse_proxy * unix/${config.services.e2clicker.frontend.socketPath}
        }
      ''
      + (lib.optionalString (noFlag "no-backend")) ''
        handle /api* {
          reverse_proxy * localhost:36001
        }
      '';
  };

  environment.systemPackages = with pkgs; [
    pgcli
    (pkgs.writeShellScriptBin "e2clicker-sql" ''
      pgcli "postgresql://root@/e2clicker-backend"
    '')
  ];

  networking.firewall.allowedTCPPorts = [ 80 ];

  programs.bash.interactiveShellInit = ''
    # Stop VM on Ctrl+D or exit
    alias __shutdown="shutdown now"
    trap  __shutdown  EXIT
  '';

  system.stateVersion = "24.11";
}
