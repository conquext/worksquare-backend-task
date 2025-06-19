package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetProjectRoot returns the project root directory
func GetProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	
	// Go up directories until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root, fallback to current directory
			return "."
		}
		dir = parent
	}
}

// GetDataFilePath returns the absolute path to a data file
func GetDataFilePath(filename string) string {
	projectRoot := GetProjectRoot()
	return filepath.Join(projectRoot, "data", filename)
}