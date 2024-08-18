{ ... }:

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
    ];
    identMap = ''
      e2clicker-backend e2clicker-backend e2clicker-backend
    '';
    extraPlugins = ps: with ps; [ ];
  };
}
