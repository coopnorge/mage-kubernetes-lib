package magekubernetes

import (
	"errors"
	"os"
)

const palletDirectory = ".pallet"

func listPalletFiles() ([]string, error) {
	_, err := os.Stat(palletDirectory)

	if errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return []string{}, err
	}
	return listFilesInDirectory(palletDirectory)
}
