package mgit

import (
	"os"
	"path/filepath"
)

func isGitRepo(directoryPath string) bool {
	gitDirPath := filepath.Join(directoryPath, ".git")
	return isDirectory(gitDirPath)
}

func fullPath(baseDirectory string, statusItem gitStatus) string {
	return filepath.Join(baseDirectory, statusItem.dirName)
}

func isDirectory(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}
