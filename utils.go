package magekubernetes

import (
	"os"
	"path/filepath"
)

func listFilesInDirectory(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func tempDir() (string, error) {
	dir, err := os.MkdirTemp("", "kubernetes-validation-*")
	if err != nil {
		return "", err
	}
	return dir, nil
}

func isKustomizeDir(dirPath string) bool {
	if _, err := os.Stat(dirPath + "/kustomization.yaml"); err == nil {
		return true
	}
	// Support legacy .yml extension
	if _, err := os.Stat(dirPath + "/kustomization.yml"); err == nil {
		return true
	}
	return false
}
