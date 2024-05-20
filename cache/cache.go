package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetCacheDir(name string) (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user's cache directory")
	}

	newCacheDir := filepath.Join(userCacheDir, name)

	err = os.MkdirAll(newCacheDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create cache directory")
	}

	return newCacheDir, nil
}
