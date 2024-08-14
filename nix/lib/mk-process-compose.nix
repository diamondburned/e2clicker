{
  pkgs,
  lib,
  self,
  ...
}:

let
  mkJSONConfig = name: config: pkgs.writeText name (builtins.toJSON config);

  processComposeFile =
    processesConfig:
    mkJSONConfig "process-compose.json" (
      let
        mapDependsOn =
          processes: condition:
          builtins.listToAttrs (map (process: lib.nameValuePair process { inherit condition; }) processes);
      in
      {
        version = "0.5";
        processes = lib.mapAttrs' (name: process: {
          inherit name;
          value = (
            {
              working_dir = "./" + (process.workingDirectory or name);
              command =
                let
                  script = pkgs.writeShellScript "${name}.sh" ''
                    set -ex
                    ${process.command}
                  '';
                in
                "nix develop -c ${script}";
              environment = lib.mapAttrsToList (k: v: "${k}=${v}") (
                (process.environment or { })
                // (lib.optionalAttrs (process ? "port") { PORT = toString process.port; })
              );
              shutdown = {
                signal = 2; # SIGINT
                timeout_seconds = 5;
              };
              depends_on =
                { }
                // (mapDependsOn (process.after or [ ]) "process_completed_successfully")
                // (mapDependsOn (process.dependsOn or [ ]) "process_healthy");
            }
            // (lib.optionalAttrs (process ? "healthPath") {
              readiness_probe = {
                http_get = {
                  host = "localhost";
                  port = process.port;
                  path = process.healthPath;
                };
                period_seconds = 1;
                failure_threshold = 300;
              };
              availability = {
                restart = "exit_on_failure";
                backoff_seconds = 10;
              };
            })
          );
        }) processesConfig;
      }
    );

  mkProcessCompose =
    processesConfig:
    pkgs.writeShellApplication rec {
      name = "twipi-dev";
      meta.mainProgram = "twipi-dev";
      runtimeInputs = [ pkgs.process-compose ];
      text = ''
        exec process-compose up \
          --config ${processComposeFile processesConfig} \
          --ref-rate 100ms \
          --no-server \
          --keep-tui
      '';
    };
in
mkProcessCompose
