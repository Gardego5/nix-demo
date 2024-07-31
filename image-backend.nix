{ pkgs, name, ... }:
let
  app = pkgs.buildGoModule {
    pname = "${name}-backend";
    src = ./backend;
    vendorHash = "sha256-SfCrX+zJDpybR3Hox/wCgMlxUThYODxt1lXd8b/6ui4=";
    version = "0.0.1";
  };
in pkgs.dockerTools.buildLayeredImage {
  name = app.pname;
  contents = with pkgs; [ bash curl ];
  config = {
    Cmd = [ "${app}/bin/backend" ];
    Env = [ "PORT=80" ];
    ExposedPorts."80/tcp" = { };
  };
}
