package mgit

import (
	"atilm/mgit/utilities"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

func extractNumberOrDefault(input string, regexPattern string, defaultValue int) int {
	numberString, err := extractString(input, regexPattern)
	if err != nil {
		return defaultValue
	}

	return toIntOrDefault(numberString, defaultValue)
}

func toIntOrDefault(input string, defaultValue int) int {
	number, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	} else {
		return number
	}
}

func extractString(input string, regexPattern string) (string, error) {
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindStringSubmatch(input)

	if len(matches) >= 2 {
		return matches[1], nil
	} else {
		return "", ErrNotFound
	}
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
			gitStatusItems = replace(gitStatusItems, status)
			printStatusItems(gitStatusItems, printer)
		case <-doneChannel:
			loop = false
		}
	}

	return nil
}

func replace(gitStatusItems []gitStatus, status gitStatus) []gitStatus {
	for i, oldStatus := range gitStatusItems {
		if oldStatus.index == status.index {
			gitStatusItems[i] = status
			break
		}
	}

	return gitStatusItems
}
