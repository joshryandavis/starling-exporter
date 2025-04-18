{
  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
  };
  outputs = inputs @ {flake-parts, ...}: let
    name = "starlingexporter";
    info = builtins.fromJSON (builtins.readFile ./nixpkg.json);
  in
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = import inputs.systems;
      imports = [inputs.flake-parts.flakeModules.easyOverlay];
      flake = {
        nixosModules.default = {
          config,
          lib,
          pkgs,
          ...
        }: let
          pkg =
            pkgs.buildGoModule
            {
              src = ./.;
              pname = name;
              version = "snapshot";
              meta.mainprogram = name;
              vendorHash = info.vendor-hash;
              hash = info.bin-hash;
            };
        in {
          options = {
            services.starlingexporter = {
              enable = lib.mkOption {
                type = lib.types.bool;
                default = false;
              };
              token = lib.mkOption {
                type = lib.types.str;
                default = "";
              };
              tokenPath = lib.mkOption {
                type = lib.types.str;
                default = "";
              };
            };
          };
          config = {
            systemd.services."starling-exporter" = {
              description = "Starling Exporter";
              wantedBy = ["multi-user.target"];
              after = ["network.target"];
              environment = {
                STARLING_ACCESS_TOKEN_PATH = config.services.starlingexporter.tokenPath;
                STARLING_ACCESS_TOKEN = config.services.starlingexporter.token;
              };
              serviceConfig = {
                Type = "simple";
                ExecStart = "${pkg}/bin/starlingexporter";
                Restart = "always";
              };
            };
          };
        };
      };
    };
}
