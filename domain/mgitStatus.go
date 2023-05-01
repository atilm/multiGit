package mgit

import (
	"atilm/mgit/utilities"
	"errors"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrNotARepository = errors.New("Not a git repository")
	ErrNotFound       = errors.New("Not found")
)

func ReportStatus(baseDirectory string, printer utilities.ConsolePrinter) error {
	gitStatusItems, err := CollectGitStatusFromSubdirectories(baseDirectory)
	if err != nil {
		return err
	}

	updateFunction := func(status gitStatus) (gitStatus, error) {
		return queryGitStatus(baseDirectory, status)
	}

	return parallelStatusUpdate(updateFunction, gitStatusItems, printer)
}

func queryGitStatus(baseDirectory string, currentStatus gitStatus) (gitStatus, error) {
	fullRepoPath := filepath.Join(baseDirectory, currentStatus.dirName)

	if err := executeFetchCommand(fullRepoPath); err != nil {
		return currentStatus, err
	}

	return executeStatusCommand(fullRepoPath, currentStatus)
}

func executeFetchCommand(repoPath string) error {
	fetchCommand := exec.Command("git", "fetch")
	fetchCommand.Dir = repoPath
	return fetchCommand.Run()
}

func executeStatusCommand(repoPath string, currentStatus gitStatus) (gitStatus, error) {
	statusCommand := exec.Command("git", "status")
	statusCommand.Dir = repoPath
	output, _ := statusCommand.CombinedOutput()
	outputString := string(output)

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
