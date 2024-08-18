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
      port = 36000;
    };
    backend = {
      enable = true;
      debug = true;
      port = 36001;
      databaseURI = "postgresql://e2clicker-backend@/e2clicker-backend";
    };
  };

  services.caddy = {
    enable = true;
    virtualHosts.":80".extraConfig = ''
      handle /api* {
        reverse_proxy * localhost:36001
      }
      handle {
        reverse_proxy * localhost:36000
      }
    '';
  };

  networking.firewall.allowedTCPPorts = [ 80 ];

  system.stateVersion = "24.11";
}
