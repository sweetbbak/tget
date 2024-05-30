package main

import (
	"fmt"
	"math"

	"github.com/anacrolix/torrent"
)

// https://gist.github.com/anikitenko/b41206a49727b83a530142c76b1cb82d
func prettyByteSize(b int) string {
	bf := float64(b)
	for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", bf)
}

func FileInfo(tor string, client *torrent.Client) error {
	t, err := AddTorrent(client, tor)
	if err != nil {
		return err
	}

	printFileInfo(t)
	t.Drop()
	client.Close()
	return nil
}

func Info(tor string, client *torrent.Client) error {
	t, err := AddTorrent(client, tor)
	if err != nil {
		return err
	}

	printInfo(t)

	t.Drop()
	client.Close()
	return nil
}

// print torrent information and files + their index
func printInfo(t *torrent.Torrent) {
	info := t.Info()
	files := t.Files()
	sz := info.TotalLength()
	psz := prettyByteSize(int(sz))
	fmt.Printf("[%s] %s\n", psz, info.Name)

	fmt.Printf("Files:\n")
	for i, f := range files {
		fmt.Printf("%03d | %s\n", i, f.DisplayPath())
	}
}

func printFileInfo(t *torrent.Torrent) {
	info := t.Info()
	files := t.Files()
	sz := info.TotalLength()
	psz := prettyByteSize(int(sz))
	fmt.Printf("%s\n", psz)

	for i, f := range files {
		fmt.Printf("%d | %s\n", i, f.DisplayPath())
	}
}
