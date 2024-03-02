{ lib, fetchFromGitHub, buildGoModule }:

buildGoModule rec {
    pname = "tget";
    version = "0.1";

    src = fetchFromGitHub {
      owner = "sweetbbak";
      repo = "tget";
      # rev = "9dc67d420208bb5f9debd260170d54035242c7ab";
      # hash = "sha256-2Z5agQtF6p21rnAcjsRr+3QOJ0QGveKVH8e9LHpm3ZE=";
    };

    vendorHash = lib.fakeHash;
    # vendorHash = "sha256-alC4/2wTbjJYWGzTDTgQweOicN3xSqfnncok/j16+0E=";

    CGO_ENABLED = 0;
    ldflags = [ "-s" "-w" ];

    tags = [ "torrent" "bittorrent" "anime" ];
    proxyVendor = true;

    buildPhase = ''
        go mod vendor
        go build
        upx -9 tget
    '';

    installPhase = ''
        mkdir -p $out/bin
        mv toru $out/bin
    '';

    meta = with lib; {
        homepage = "https://github.com/sweetbbak/tget";
        description = "wget but for torrents";
        license = licenses.mit;
        maintainers = with maintainers; [ sweetbbak ];
        mainProgram = "toru";
    };
}
