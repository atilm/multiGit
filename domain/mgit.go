package mgit

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNotARepository = errors.New("Not a git repository")
	ErrNotFound       = errors.New("Not found")
)

type gitStatus struct {
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

	if s.commitsToPull == 0 && s.commitsToPush == 0 {
		return fmt.Sprintf("%s (%s) %s[ok]",
			s.dirName,
			s.branchName,
			localChangesIndicator)
	} else {
		return fmt.Sprintf("%s (%s) %s[%d to pull, %d to push]",
			s.dirName,
			s.branchName,
			localChangesIndicator,
			s.commitsToPull,
			s.commitsToPush)
	}
}

func ReportStatus(baseDirectory string) (string, error) {
	dirEntries, err := os.ReadDir(baseDirectory)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	var repoIndex int = 0

	// statusChannel := make(chan []gitStatus)
	// doneChannel := make(chan []struct{})
	// wg := sync.WaitGroup{}

	for _, entry := range dirEntries {
		fullPath := filepath.Join(baseDirectory, entry.Name())
		if entry.IsDir() && isGitRepo(fullPath) {
			repoIndex++
			gitStatus, _ := queryGitStatus(baseDirectory, entry.Name())
			buffer.WriteString(fmt.Sprintf("%02d: %v\n", repoIndex, gitStatus.String()))
		}
	}

	return buffer.String(), nil
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

func queryGitStatus(baseDirectory string, directory string) (gitStatus, error) {
	fullRepoPath := filepath.Join(baseDirectory, directory)

	fetchCommand := exec.Command("git", "fetch")
	fetchCommand.Dir = fullRepoPath
	fetchCommand.Run()

	statusCommand := exec.Command("git", "status")
	statusCommand.Dir = fullRepoPath
	output, _ := statusCommand.CombinedOutput()
	outputString := string(output)

	// fmt.Print(outputString)

	if strings.Contains(outputString, "fatal: not a git repository") {
		return gitStatus{}, ErrNotARepository
	} else {
		branchName := extractBranchName(outputString)
		commitsToPull, commitsToPush := extractChanges(outputString)
		localChanges := hasLocalChanges(outputString)
		return gitStatus{directory, branchName, localChanges, commitsToPull, commitsToPush}, nil
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
