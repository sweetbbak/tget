{ lib, fetchFromGitHub, buildGoModule, upx }:

buildGoModule rec {
    pname = "tget";
    version = "0.1";

    src = fetchFromGitHub {
      owner = "sweetbbak";
      repo = "tget";
      rev = "959e5f4c89156b789e35a92438735004716d12a9";
      hash = "sha256-3pf5woxqUf7RfX32P21UwOBpsW1i6nkYdD6We0YjdFQ=";
    };

    # vendorHash = lib.fakeHash;
    vendorHash = "sha256-74++inwJPbpjPrK5Xn66t+s50wbA2H1RgCxrS7DVJiA=";

    CGO_ENABLED = 0;
    ldflags = [ "-s" "-w" ];
    # buildInputs = [ "upx" ];

    tags = [ "torrent" "bittorrent" "anime" ];
    proxyVendor = true;

    buildPhase = ''
        go mod vendor
        go build
    '';

    installPhase = ''
        ${upx}/bin/upx -9 tget
        mkdir -p $out/bin
        mv tget $out/bin
    '';

    meta = with lib; {
        homepage = "https://github.com/sweetbbak/tget";
        description = "wget but for torrents";
        license = licenses.mit;
        maintainers = with maintainers; [ sweetbbak ];
        mainProgram = "tget";
    };
}
