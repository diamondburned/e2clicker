{ pkgs }:

let
  flags = [
    "--no-comma-break"
    "--nogrouping"
    "--keyword-case=2"
    "--type-case=1"
    "--function-case=1"
    "--no-extra-line"
    "--spaces=2"
    "--wrap-limit=100"
    "--wrap-after=20"
    "--wrap-comment"
    "--extra-function=${
      pkgs.writeText "extra-sql-functions.txt" (
        builtins.concatStringsSep " " [
          "sqlc.embed"
        ]
      )
    }"
  ];

  flags' = lib.escapeShellArgs flags;

  inherit (pkgs) lib pgformatter;
in

pkgs.writeShellScriptBin "pg_format" ''
  exec ${lib.getExe pgformatter} ${flags'} "$@" 
''
