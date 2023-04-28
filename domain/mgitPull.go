package mgit

import (
	"atilm/mgit/utilities"
	"fmt"
	"os/exec"
	"path/filepath"
)

func Pull(baseDirectory string, args []string, printer utilities.ConsolePrinter) error {
	gitStatusItems, err := initializeStatusSlice(baseDirectory)
	if err != nil {
		return err
	}

	for _, statusItem := range gitStatusItems {
		err := executeGitPull(baseDirectory, statusItem)
		if err != nil {
			return err
		}
	}

	lines := make([]string, 0, len(gitStatusItems))
	for _, statusItem := range gitStatusItems {
		lines = append(lines, fmt.Sprintf("%02d: %s [done]",
			statusItem.index+1,
			statusItem.dirName))
	}

	printer.PrintLines(lines)

	return nil
}

func executeGitPull(baseDirectory string, statusItem gitStatus) error {
	fullRepoPath := filepath.Join(baseDirectory, statusItem.dirName)

	pullCommand := exec.Command("git", "pull")
	pullCommand.Dir = fullRepoPath
	pullCommand.Run()

	return nil
}
