package utils

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

var paths []string

func init() {
	appHomeDir, err := GetLocalDir()
	if err != nil {
		appHomeDir = ""
	}

	sharedpath, err := filepath.Abs("shared")
	if err != nil {
		sharedpath = ""
	}
	paths = []string{
		sharedpath,
		appHomeDir,
		"/etc/vhost",
	}
}

// GetLocalDir returns apps directory in user's home.
func GetLocalDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.vhost", user.HomeDir), nil
}

// DirectoryExists as path
func DirectoryExists(path string) bool {
	var stat os.FileInfo
	var err error
	if stat, err = os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return stat.IsDir()
}

// FileExists as path
func FileExists(path string) bool {
	var stat os.FileInfo
	var err error
	if stat, err = os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}

// FindFile in paths
func FindFile(filename string, paths []string) (string, error) {
	for _, filepath := range paths {
		fullPath := path.Join(filepath, filename)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("%s not found in paths: %s", filename, paths)
}

// FindFileInApp looks for files in application directories
func FindFileInApp(filename string) (string, error) {
	return FindFile(filename, paths)
}
