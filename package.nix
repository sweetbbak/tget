{
  lib,
  fetchFromGitHub,
  buildGoModule,
  upx,
}:
buildGoModule rec {
  pname = "tget";
  version = "0.1";

  src = fetchFromGitHub {
    owner = "sweetbbak";
    repo = "tget";
    rev = "9605e994b4926a8cd001e5b78c91fe7871552e5f";
    # hash = lib.fakeHash;
    hash = "sha256-gI3I/Jt8E2URvGyoo23HyWhg/qbnLvF8OsmJyeyA2xE=";
  };

  # vendorHash = lib.fakeHash;
  vendorHash = "sha256-74++inwJPbpjPrK5Xn66t+s50wbA2H1RgCxrS7DVJiA=";

  CGO_ENABLED = 0;
  ldflags = ["-s" "-w"];

  tags = ["torrent" "bittorrent" "anime"];
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
    maintainers = with maintainers; [sweetbbak];
    mainProgram = "tget";
  };
}
