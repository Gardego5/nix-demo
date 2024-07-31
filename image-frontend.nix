{ pkgs, name, ... }:
let
  src = ./frontend;
  pname = "${name}-frontend";
  version = "0.0.1";
  pnpmDeps = pkgs.pnpm.fetchDeps {
    inherit pname src version;
    hash = "sha256-RwQH8NUJUs5OvPNkah/kuRDCViqEeGhvq0JS78yZSok=";
  };
  nodeModules = pkgs.stdenv.mkDerivation {
    inherit pnpmDeps src version;
    nativeBuildInputs = [ pkgs.pnpm.configHook ];
    pname = "${pname}-prodDeps";
    NODE_ENV = "production";
    installPhase = "mkdir -p $out; cp -r node_modules $out";
  };
  app = pkgs.stdenv.mkDerivation {
    inherit pname pnpmDeps src version;
    nativeBuildInputs = with pkgs; [ nodejs pnpm.configHook ];
    buildPhase = "pnpm build";
    noCheck = true;
    installPhase = ''
      runHook preInstall
      mkdir -p $out/node_modules
      cp -r ${nodeModules}/node_modules build package.json public $out
      runHook postInstall
    '';
  };
in pkgs.dockerTools.buildLayeredImage {
  name = app.pname;
  contents = with pkgs; [ bash nodejs toybox ];
  config = {
    Cmd = [ "npm" "run" "start" ];
    Env = [ "NODE_ENV=production" "PORT=80" ];
    ExposedPorts = { "80/tcp" = { }; };
    WorkingDir = builtins.toString app;
  };
}
