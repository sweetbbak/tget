{
  lib,
  fetchFromGitHub,
  buildGoModule,
  just,
  upx,
}:
buildGoModule {
  pname = "tget";
  version = "0.1.1";

  src = fetchFromGitHub {
    owner = "sweetbbak";
    repo = "tget";
    rev = "135167df0afd3fc5df92e44b94401c223d7ee0ac";
    # hash = lib.fakeHash;
    hash = "sha256-XbIw1wLJ3PdhkOQca8UzZ1livauTdTdZbTWXrA2Nrfk=";
  };

  # vendorHash = lib.fakeHash;
  vendorHash = "sha256-2lzMwp0XN7pC5s2PYqyN+BUeCqDsFu/sMmcUorr71BY=";

  CGO_ENABLED = 0;
  ldflags = ["-s" "-w"];

  tags = ["torrent" "bittorrent" "anime"];
  proxyVendor = true;

  buildPhase = ''
    go mod vendor
    ${just}/bin/just
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
