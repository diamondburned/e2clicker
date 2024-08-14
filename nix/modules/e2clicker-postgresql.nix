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
    ensureDatabases = [ "e2clicker" ];
    ensureUsers = [
      {
        name = "e2clicker";
        ensureDBOwnership = true;
      }
    ];
    identMap = ''
      e2clicker e2clicker e2clicker
    '';
    extraPlugins = ps: with ps; [ ];
  };
}
