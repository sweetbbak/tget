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
    rev = "0a7cc2fd9eeb6ffc90f0687f580b786e76a7a90d";
    # hash = lib.fakeHash;
    hash = "sha256-9opGtMlQt3qeW/D/mVrbnBFBsGeNGk8F0YexVlLs0wE=";
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
