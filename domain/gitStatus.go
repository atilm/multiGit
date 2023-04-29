package mgit

import (
	"atilm/mgit/utilities"
	"fmt"
	"os"
	"path/filepath"
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

func CollectGitStatusFromSubdirectories(baseDirectory string) ([]gitStatus, error) {
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
