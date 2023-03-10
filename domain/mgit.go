package mgit

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrNotARepository = errors.New("Not a git repository")
)

type gitStatus struct {
	dirName    string
	branchName string
}

func queryGitStatus(baseDirectory string, directory string) (gitStatus, error) {
	statusCommand := exec.Command("git", "status")
	statusCommand.Dir = filepath.Join(baseDirectory, directory)
	output, _ := statusCommand.CombinedOutput()
	outputString := string(output)

	if strings.Contains(outputString, "fatal: not a git repository") {
		return gitStatus{}, ErrNotARepository
	} else {
		branchName := extractBranchName(outputString)
		return gitStatus{directory, branchName}, nil
	}
}

func extractBranchName(gitStatusOutput string) string {
	branchRegex := regexp.MustCompile("(?m)On branch (.+)$")
	branchNames := branchRegex.FindStringSubmatch(gitStatusOutput)

	var branchName string = "unknwon"
	if len(branchNames) >= 2 {
		branchName = branchNames[1]
	}
	return branchName
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
				buffer.WriteString(fmt.Sprintf("%02d: %s (%s) [ok]\n", repoIndex, gitStatus.dirName, gitStatus.branchName))
			}
		}
	}

	return buffer.String(), nil
}
