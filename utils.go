package magekubernetes

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func listFilesInDirectory(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, _ error) error {
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

func getBoolEnv(key string, defaultValue bool) (bool, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultValue, err
	}
	return b, nil
}

func runInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
