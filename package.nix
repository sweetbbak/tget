{ lib, fetchFromGitHub, buildGoModule }:

buildGoModule rec {
    pname = "tget";
    version = "0.1";

    src = fetchFromGitHub {
      owner = "sweetbbak";
      repo = "tget";
      rev = "48ed9a95fc68c9455f6f063dec387b1f2bfa441f";
      hash = "sha256-3pf5woxqUf7RfX32P21UwOBpsW1i6nkYdD6We0YjdFQ=";
      # hash = "sha256-2Z5agQtF6p21rnAcjsRr+3QOJ0QGveKVH8e9LHpm3ZE=";
    };

    # vendorHash = lib.fakeHash;
    vendorHash = "sha256-74++inwJPbpjPrK5Xn66t+s50wbA2H1RgCxrS7DVJiA=";

    CGO_ENABLED = 0;
    ldflags = [ "-s" "-w" ];

    tags = [ "torrent" "bittorrent" "anime" ];
    proxyVendor = true;

    buildPhase = ''
        go mod vendor
        go build
    '';

    installPhase = ''
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
