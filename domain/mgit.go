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
		return fmt.Sprintf("%s (%s) [%d to pull, %d to push]",
			s.dirName,
			s.branchName,
			s.commitsToPull,
			s.commitsToPush)
	}
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
		localChanges := strings.Contains(outputString, "Untracked files:")
		return gitStatus{directory, branchName, localChanges, commitsToPull, commitsToPush}, nil
	}
}

func extractBranchName(gitStatusOutput string) string {
	branchName, err := extractString(gitStatusOutput, "(?m)On branch (.+)$")

	if err == nil {
		return branchName
	}

	return "unknown"
}

func extractChanges(gitStatusOutput string) (int, int) {
	numberString, err := extractString(gitStatusOutput, `Your branch is behind .+ by (\d+) commit`)

	if err == nil {
		commitsToPull, err := strconv.Atoi(numberString)
		if err != nil {
			return 0, 0
		} else {
			return commitsToPull, 0
		}
	} else {
		return 0, 0
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

func ReportStatus(baseDirectory string) (string, error) {
	dirEntries, err := os.ReadDir(baseDirectory)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	var repoIndex int = 0

	for _, entry := range dirEntries {
		if entry.IsDir() {
			gitStatus, err := queryGitStatus(baseDirectory, entry.Name())
			if err == nil {
				repoIndex++
				buffer.WriteString(fmt.Sprintf("%02d: %v\n", repoIndex, gitStatus.String()))
			}
		}
	}

	return buffer.String(), nil
}
