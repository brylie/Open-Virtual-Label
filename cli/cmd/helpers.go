package cmd

import "os"

// statFile is a thin wrapper around os.Stat for use in cmd files.
func statFile(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// osCreateDir creates a directory and all parents.
func osCreateDir(path string) error {
	return os.MkdirAll(path, 0o755)
}
