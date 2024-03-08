package main

import (
	"io"
	"os"
	"strings"

	"github.com/anacrolix/torrent/iplist"
)

func parseBlocklist(r io.Reader) (*iplist.IPList, error) {
	i, err := iplist.NewFromReader(r)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func openBlocklist(blist string) (io.Reader, error) {
	var fpath string
	var err error

	if strings.HasPrefix(blist, "http") {
		fpath, err = fileFromURL(blist)
		if err != nil {
			return nil, err
		}
	} else {
		fpath = blist
	}

	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	return f, nil
}
