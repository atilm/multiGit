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
)

type gitStatus struct {
	dirName       string
	branchName    string
	localChanges  bool
	commitsToPull int
	commitsToPush int
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
	branchRegex := regexp.MustCompile("(?m)On branch (.+)$")
	branchNames := branchRegex.FindStringSubmatch(gitStatusOutput)

	var branchName string = "unknown"
	if len(branchNames) >= 2 {
		branchName = branchNames[1]
	}
	return branchName
}

func extractChanges(gitStatusOutput string) (int, int) {
	behindRegex := regexp.MustCompile(`Your branch is behind .+ by (\d+) commit`)
	numbers := behindRegex.FindStringSubmatch(gitStatusOutput)

	if len(numbers) >= 2 {
		commitsToPull, err := strconv.Atoi(numbers[1])
		if err != nil {
			return 0, 0
		} else {
			return commitsToPull, 0
		}
	} else {
		return 0, 0
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

				var localChangesIndicator string = ""
				if gitStatus.localChanges {
					localChangesIndicator = "* "
				}

				if gitStatus.commitsToPull == 0 && gitStatus.commitsToPush == 0 {
					buffer.WriteString(fmt.Sprintf("%02d: %s (%s) %s[ok]\n", repoIndex, gitStatus.dirName, gitStatus.branchName, localChangesIndicator))
				} else {
					buffer.WriteString(fmt.Sprintf("%02d: %s (%s) [%d to pull, %d to push]\n", repoIndex, gitStatus.dirName, gitStatus.branchName,
						gitStatus.commitsToPull,
						gitStatus.commitsToPush))
				}
			}
		}
	}

	return buffer.String(), nil
}
