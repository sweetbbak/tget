package main

import (
	"log"

	"github.com/anacrolix/torrent/storage"
)

// getStorage returns a storage implementation that writes downloaded
// files to a user-defined directory, and writes metadata files to a
// temporary directory.
func getMetadataDir(metadataDir, downloadDir string) (storage.ClientImpl, error) {
	mstor, err := storage.NewDefaultPieceCompletionForDir(metadataDir)
	if err != nil {
		log.Println(err)
		return storage.NewMMap(downloadDir), nil
	}

	tstor := storage.NewMMapWithCompletion(downloadDir, mstor)
	if err != nil {
		return nil, err
	}

	return tstor, err
}
