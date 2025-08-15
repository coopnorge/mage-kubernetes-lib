package magekubernetes

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/magefile/mage/sh"
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

func since(msg string, start time.Time, fields map[string]any) {
	d := time.Since(start)
	var parts []string
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	if len(parts) > 0 {
		infof("%s done in %s (%s)", msg, d, strings.Join(parts, " "))
	} else {
		infof("%s done in %s", msg, d)
	}
}

// runLogged logs the exact command (with safe shell-style quoting) and runs it.
func runLogged(name string, args ...string) error {
	cmdline := strings.Join(append([]string{name}, args...), " ")
	debugf("Running: %s", cmdline)
	if err := sh.Run(name, args...); err != nil {
		errorf("Command failed: %s: %v", cmdline, err)
		return err
	}
	return nil
}
