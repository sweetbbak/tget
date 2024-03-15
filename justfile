default:
    go build -ldflags "-s -w" ./cmd/tget
pack:
    go build -ldflags "-s -w" ./cmd/tget
    upx -9 tget
