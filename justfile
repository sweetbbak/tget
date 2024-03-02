default:
    go build -ldflags "-s -w" .
pack:
    go build -ldflags "-s -w" .
    upx -9 tget
