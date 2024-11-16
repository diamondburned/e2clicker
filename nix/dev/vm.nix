{
  config,
  lib,
  pkgs,
  ...
}@inputs:

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
      enable = true;
      socket = true;
    };
    backend = {
      enable = true;
      debug = true;
      api = {
        listenAddress = ":36001";
      };
      postgresql = {
        databaseURI = "postgresql://e2clicker-backend@/e2clicker-backend";
      };
      notification = {
        clientTimeout = "15s";
        webPushKeys =
          let
            path = ../../vapid-keys.json;
          in
          if builtins.pathExists path then path else null;
      };
    };
  };

  services.caddy = {
    enable = true;
    virtualHosts.":80".extraConfig = ''
      handle /_app/immutable* {
        header Cache-Control "public, immutable, max-age=31536000"
        file_server {
          root ${config.services.e2clicker.frontend.package.assets}
          precompressed gzip br
          pass_thru
        }
      }
      handle /api* {
        reverse_proxy * localhost:36001
      }
      handle {
        reverse_proxy * unix/${config.services.e2clicker.frontend.socketPath}
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
