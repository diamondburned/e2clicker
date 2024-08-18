{ pkgs, self, ... }:

let
  lib = pkgs.lib;
in

pkgs.writeShellScriptBin "e2clicker-dev" ''
  ${lib.getExe pkgs.zellij} delete-session e2clicker
  ${lib.getExe pkgs.zellij} \
    --layout ${./zellij-layout.kdl} \
    --session e2clicker
''
