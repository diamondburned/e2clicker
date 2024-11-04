{
  config,
  lib,
  pkgs,
  ...
}:

{
  services.postgresql = {
    enable = true;
    ensureDatabases = [ "e2clicker-backend" ];
    ensureUsers = [
      {
        name = "e2clicker-backend";
        ensureDBOwnership = true;
      }
      {
        name = "root";
        ensureClauses = {
          login = true;
          superuser = true;
        };
      }
    ];
    identMap = ''
      e2clicker-backend e2clicker-backend e2clicker-backend
    '';
    extraPlugins = ps: with ps; [ ];
  };

  environment.systemPackages = with pkgs; [ pgcli ];
}
