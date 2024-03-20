default:
    CGO_ENABLED=0 go build -ldflags "-s -w" ./cmd/tget

pack:
    go build -ldflags "-s -w" ./cmd/tget
    upx -9 tget
