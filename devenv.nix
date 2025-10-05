{
  config,
  lib,
  pkgs,
  ...
}:
let
  appBinaryPath = "./3d-rack-brackets";
  outputScadPath = "./output/output.scad";
  previewPngPath = "./output.png";
  outputStlPath = "./output/output.stl";
in
{
  packages = with pkgs; [
    golangci-lint
    graphviz
    inotify-tools
  ];

  languages.go = {
    enable = true;
    enableHardeningWorkaround = true;
  };

  processes =
    with pkgs;
    lib.optionalAttrs (!config.devenv.isTesting) {
      docs.exec = "${lib.getExe go} doc -http";

      open = {
        exec = "${lib.getExe openscad} ${outputScadPath}";
      };

      profile = {
        exec = "${appBinaryPath} --cpu-profile=./output/render.prof render --production ${outputScadPath} && ${lib.getExe go} tool pprof -http 0.0.0.0:8080 ${appBinaryPath} ./output/render.prof";
      };

      watch = {
        exec = "./build/render-loop.bash";
      };
    };

  tasks =
    let
      goInputFiles = [
        "main.go"
        "go.mod"
        "go.sum"
        "./internal/**/*.go"
      ];
    in
    with pkgs;
    {
      "devenv:processes:open" = {
        after = [ "app:render" ];
      };
      "devenv:processes:profile" = {
        after = [ "app:build" ];
      };

      "app:build".exec = "${lib.getExe go} build .";

      "app:lint".exec = "${lib.getExe golangci-lint} run ./...";
      "app:lint-fix".exec = "${lib.getExe golangci-lint} run --fix ./...";

      "app:makeOutputDir".exec = "mkdir -p output";

      "app:render" = {
        exec = "${lib.getExe go} run . render ${outputScadPath}";
        after = [ "app:makeOutputDir" ];
        execIfModified = goInputFiles;
      };
      "app:render-debug" = {
        exec = "${lib.getExe go} run . --debug render ${outputScadPath}";
        after = [ "app:makeOutputDir" ];
      };
      "app:render-prod" = {
        exec = "${lib.getExe go} run . render --production ${outputScadPath}";
        after = [ "app:makeOutputDir" ];
      };

      "app:render-preview" = {
        exec = "${lib.getExe openscad} -o ${previewPngPath} ${outputScadPath}";
        after = [ "app:render-prod" ];
      };
      "app:render-stl" = {
        exec = "${lib.getExe openscad} -o ${outputStlPath} ${outputScadPath}";
        after = [ "app:render-prod" ];
      };

      "app:test" = {
        exec = "${lib.getExe go} test ./...";
      };
    };

  git-hooks.hooks = {
    check-merge-conflicts.enable = true;
    check-shebang-scripts-are-executable.enable = true;
    end-of-file-fixer.enable = true;
    gitleaks = {
      enable = true;
      name = "Gitleaks";
      entry = "${lib.getExe pkgs.gitleaks} git --pre-commit --redact --staged --verbose";
      pass_filenames = false;
    };
    gofmt.enable = true;
    golangci-lint.enable = true;
    nixfmt-rfc-style.enable = true;
    shellcheck.enable = true;
  };

  outputs = {
    "3d-rack-brackets" = pkgs.buildGo124Module {
      name = "3d-rack-brackets";
      src = builtins.path {
        path = ./.;
        name = "source";
      };
      vendorHash = "sha256-OX8z1bp3Kzsz6phWnl8EWPgwiH2VJlHkFWt1qzfwTOg=";
    };
  };
}
