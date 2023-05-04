package util

import (
	"os"
	"path/filepath"
	"strings"
)

func AbsPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") || path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path, err
		}
		path = filepath.Join(homeDir, path[2:])
	}
	return path, nil
}
