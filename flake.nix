{
  description = "Convert BEP to a useful Parquet file with Go tools";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [
          (final: prev: {
            build-event-protocol-analysis-tools = final.callPackage ./derivation.nix {};
          })
        ];
      };
    in {
      formatter = pkgs.alejandra;

      packages.default = pkgs.build-event-protocol-analysis-tools;

      devShells.default = pkgs.mkShell {
        inputsFrom = with pkgs; [build-event-protocol-analysis-tools];
        buildInputs = with pkgs; [
          gopls
          gotools
          duckdb
        ];
        shellHook = ''
          echo "[+] Generating Protobufs..."
          mkdir -p genproto
          protoc -I=proto/ \
                 --go_opt=module=github.com/fzakaria/build-event-protocol-analysis-tools/genproto \
                 --go_out=genproto proto/*.proto
          echo "[+] Done."
        '';
      };
    });
}
