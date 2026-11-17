{
  pkgs,
  self,
  inputs,
  ...
}:

let
  testVapidKeys = pkgs.writeText "e2clicker-test-vapid.json" (
    builtins.toJSON {
      publicKey = "BH1uVr_j2jFqnZFmOaqvBgwLZhMIRgGf23u-DfrN8yyr_0WX4UsgVIhWil6z__ddipWL4E8LqGgAc9yxIuTaJz8";
      privateKey = "YqR2fQc-TRTa40lhrjJiLJFNGtPCLEekn2fgYZxZQAA";
      subscriber = "mailto:test@e2clicker.local";
    }
  );
in
pkgs.testers.runNixOSTest {
  name = "e2clicker-e2e";

  node.specialArgs = {
    inherit self inputs;
  };

  nodes.machine =
    { config, pkgs, ... }:
    {
      imports = [
        self.nixosModules.e2clicker
        self.nixosModules.e2clicker-postgresql
      ];

      services.e2clicker = {
        frontend = {
          enable = true;
          socket = true;
          trustProxy = true;
        };
        backend = {
          enable = true;
          debug = true;
          logFormat = "text";
          api.listenAddress = "127.0.0.1:36001";
          postgresql.databaseURI = "postgresql://e2clicker-backend@/e2clicker-backend";
          notification = {
            clientTimeout = "5s";
            webPush = "${testVapidKeys}";
          };
        };
      };

      services.caddy = {
        enable = true;
        globalConfig = "auto_https off";
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
            reverse_proxy * 127.0.0.1:36001
          }
          handle {
            reverse_proxy * unix/${config.services.e2clicker.frontend.socketPath}
          }
        '';
      };

      environment.etc."e2clicker-tests/api".source = ../api;

      environment.systemPackages = [
        pkgs.curl
        pkgs.hurl
        config.services.postgresql.package
      ];

      documentation.enable = false;
      virtualisation.memorySize = 1024;
    };

  testScript = builtins.readFile ./test_script.py;
}
