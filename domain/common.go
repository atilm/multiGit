package mgit

import (
	"atilm/mgit/utilities"
	"os"
	"path/filepath"
	"sync"
)

func fullPath(baseDirectory string, statusItem gitStatus) string {
	return filepath.Join(baseDirectory, statusItem.dirName)
}

func isDirectory(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

func isGitRepo(directoryPath string) bool {
	gitDirPath := filepath.Join(directoryPath, ".git")
	return isDirectory(gitDirPath)
}

func parallelStatusUpdate(
	updateFunc func(current gitStatus) (gitStatus, error),
	gitStatusItems []gitStatus,
	printer utilities.ConsolePrinter) error {

	statusChannel := make(chan gitStatus)
	doneChannel := make(chan struct{})
	wg := sync.WaitGroup{}

	for _, item := range gitStatusItems {
		wg.Add(1)
		go func(currentStatus gitStatus) {
			defer wg.Done()
			status, _ := updateFunc(currentStatus)
			statusChannel <- status
		}(item)
	}

	go func() {
		wg.Wait()
		close(doneChannel)
	}()

	printStatusItems(gitStatusItems, printer)
	loop := true
	for loop {
		select {
		case status := <-statusChannel:
			gitStatusItems[status.index] = status
			printStatusItems(gitStatusItems, printer)
		case <-doneChannel:
			loop = false
		}
	}

	return nil
}
