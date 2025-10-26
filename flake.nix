{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    let
      pkgs-amd64 = nixpkgs.legacyPackages.x86_64-linux;
      pkgs-mister = nixpkgs.legacyPackages.x86_64-linux.pkgsCross.armv7l-hf-multiplatform;
    in
    {
      packages =
        let
          version = "1.2.3";
        in
        rec {
          x86_64-linux = rec {
            default = pkgs-amd64.buildGoModule {
              name = "mister-mqtt";
              src = ./.;
              vendorHash = null;
            };
            release = pkgs-amd64.runCommand "release" { } ''
              mkdir $out
              cp "${downloader}" "$out/mister-mqtt.json"
              cp "${armv7l-linux.default}/bin/mister-mqtt" "$out/mister-mqtt"
            '';
            downloader = pkgs-amd64.writeText "mister-mqtt.json" (
              builtins.toJSON {
                db_id = "ChrisOboe/mister-mqtt";
                timestamp = self.lastModified;
                folders."daemons/" = { };
                files."daemons/mister-mqtt" = {
                  hash = builtins.readFile (
                    pkgs-amd64.runCommand ''hash'' { }
                      ''${pkgs-amd64.coreutils}/bin/md5sum "${armv7l-linux.default}/bin/mister-mqtt" | ${pkgs-amd64.coreutils}/bin/cut -c -32> $out''
                  );
                  size = builtins.fromJSON (
                    builtins.readFile (
                      pkgs-amd64.runCommand ''size'' { }
                        ''${pkgs-amd64.coreutils}/bin/stat --printf="%s" "${armv7l-linux.default}/bin/mister-mqtt" > $out''
                    )
                  );
                  url = "https://github.com/ChrisOboe/mister-mqtt/releases/download/${version}/mister-mqtt";
                  reboot = true;
                  tags = [ "mister-mqtt" ];
                };
              }
            );
          };
          armv7l-linux.default = pkgs-mister.buildGoModule {
            name = "mister-mqtt";
            src = ./.;
            vendorHash = null;
          };
        };

      devShells.x86_64-linux.default = pkgs-amd64.mkShell {
        buildInputs = with pkgs-amd64; [
          go
          gopls
          go-tools
          golangci-lint
          git
        ];

        shellHook = ''
          echo "Welcome to mister-mqtt development environment!"
          echo "Go version: $(go version)"
        '';
      };
    };
}
