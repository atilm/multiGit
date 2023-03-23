package mgit

import (
	"atilm/mgit/utilities"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrNotARepository = errors.New("Not a git repository")
	ErrNotFound       = errors.New("Not found")
)

type gitStatus struct {
	index         uint
	dirName       string
	branchName    string
	localChanges  bool
	commitsToPull int
	commitsToPush int
}

func (s *gitStatus) String() string {
	var localChangesIndicator string = ""

	if s.localChanges {
		localChangesIndicator = "* "
	}

	entryId := fmt.Sprintf("%02d: %s",
		s.index+1,
		s.dirName)

	if s.IsUninitialized() {
		return fmt.Sprintf("%s ...", entryId)
	}

	entryBegin := fmt.Sprintf("%s (%s) %s",
		entryId,
		s.branchName,
		localChangesIndicator)

	if s.IsInSync() {
		return entryBegin + "[ok]"
	}

	return entryBegin + fmt.Sprintf("[%d to pull, %d to push]",
		s.commitsToPull,
		s.commitsToPush)
}

func (s *gitStatus) IsUninitialized() bool {
	return s.branchName == ""
}

func (s *gitStatus) IsInSync() bool {
	return s.commitsToPull == 0 && s.commitsToPush == 0
}

func ReportStatus(baseDirectory string, printer utilities.ConsolePrinter) error {
	gitStatusItems, err := initializeStatusSlice(baseDirectory)
	if err != nil {
		return err
	}

	statusChannel := make(chan gitStatus)
	doneChannel := make(chan struct{})
	wg := sync.WaitGroup{}

	for _, item := range gitStatusItems {
		wg.Add(1)
		go func(currentStatus gitStatus) {
			defer wg.Done()
			status, _ := queryGitStatus(baseDirectory, currentStatus)
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

func printStatusItems(items []gitStatus, printer utilities.ConsolePrinter) {
	lines := createStatusLines(items)
	printer.PrintLines(lines)
}

func createStatusLines(items []gitStatus) []string {
	lines := make([]string, 0, len(items))
	for _, status := range items {
		lines = append(lines, status.String())
	}
	return lines
}

func initializeStatusSlice(baseDirectory string) ([]gitStatus, error) {
	statusEntries := make([]gitStatus, 0, 10)

	dirEntries, err := os.ReadDir(baseDirectory)
	if err != nil {
		return statusEntries, err
	}

	var repoIndex int = 0
	for _, entry := range dirEntries {
		fullPath := filepath.Join(baseDirectory, entry.Name())
		if entry.IsDir() && isGitRepo(fullPath) {
			statusEntries = append(statusEntries, gitStatus{index: uint(repoIndex), dirName: entry.Name()})
			repoIndex++
		}
	}

	return statusEntries, nil
}

func isGitRepo(directoryPath string) bool {
	gitDirPath := filepath.Join(directoryPath, ".git")
	return isDirectory(gitDirPath)
}

func isDirectory(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

func queryGitStatus(baseDirectory string, currentStatus gitStatus) (gitStatus, error) {
	fullRepoPath := filepath.Join(baseDirectory, currentStatus.dirName)

	fetchCommand := exec.Command("git", "fetch")
	fetchCommand.Dir = fullRepoPath
	fetchCommand.Run()

	statusCommand := exec.Command("git", "status")
	statusCommand.Dir = fullRepoPath
	output, _ := statusCommand.CombinedOutput()
	outputString := string(output)

	// fmt.Print(outputString)

	if strings.Contains(outputString, "fatal: not a git repository") {
		return currentStatus, ErrNotARepository
	} else {
		branchName := extractBranchName(outputString)
		commitsToPull, commitsToPush := extractChanges(outputString)
		localChanges := hasLocalChanges(outputString)
		return gitStatus{currentStatus.index, currentStatus.dirName, branchName, localChanges, commitsToPull, commitsToPush}, nil
	}
}

func hasLocalChanges(gitStatusOutput string) bool {
	return strings.Contains(gitStatusOutput, "Untracked files:") ||
		strings.Contains(gitStatusOutput, "Changes not staged for commit:") ||
		strings.Contains(gitStatusOutput, "Changes to be committed:")
}

func extractBranchName(gitStatusOutput string) string {
	branchName, err := extractString(gitStatusOutput, "(?m)On branch (.+)$")

	if err == nil {
		return branchName
	}

	return "unknown"
}

func extractChanges(gitStatusOutput string) (int, int) {
	// Your branch and 'origin/main' have diverged, ...
	regex := regexp.MustCompile(`and have (\d+) and (\d+) different commits each`)
	matches := regex.FindStringSubmatch(gitStatusOutput)
	if len(matches) >= 3 {
		return toIntOrDefault(matches[1], 0), toIntOrDefault(matches[2], 0)
	}

	commitsToPull := extractNumberOrDefault(gitStatusOutput, `Your branch is behind .+ by (\d+) commit`, 0)
	commitsToPush := extractNumberOrDefault(gitStatusOutput, `Your branch is ahead of .+ by (\d+) commit`, 0)

	return commitsToPull, commitsToPush
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
