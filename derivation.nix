{
  buildGoModule,
  lib,
  protobuf,
  protoc-gen-go,
}: let
  fs = lib.fileset;
in
  buildGoModule {
    name = "build-event-protocol-analysis-tools";
    version = "0.0.1";
    src = fs.toSource {
      root = ./.;
      fileset = fs.unions [
        ./proto
        ./cmd
        ./go.mod
        ./go.sum
      ];
    };
    vendorHash = "sha256-i+q1dgte6rJGdIv6DsRXmKA3vM8Vj9wB8BCluqaFTnI=";
    nativeBuildInputs = [
      protobuf
      protoc-gen-go
    ];
    preBuild = ''
      mkdir genproto
      protoc -I=proto/ --go_opt=module=github.com/fzakaria/build-event-protocol-analysis-tools/genproto \
             --go_out=genproto proto/*.proto
    '';
  }
