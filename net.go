package main

import (
	"io"
	"net/http"
	"os"
	"path"
)

func fileFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fname := path.Base(url)
	tmp := os.TempDir()
	fname = path.Join(tmp, fname)

	file, err := os.Create(fname)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	} else {
		return fname, nil
	}
}
