let
  port = 8000;
in

{
  containers.e2clicker = {
    config =
      {
        config,
        lib,
        pkgs,
        self,
        inputs,
        ...
      }:
      {
        imports = [
          "${self}/nix/modules/e2clicker.nix"
          "${self}/nix/modules/postgresql.nix"
        ];

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
            enable = false;
          };
          backend = {
            enable = true;
            httpAddress = "localhost:36001";
            databaseURI = "postgresql://e2clicker@localhost";
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
      };
  };
}
