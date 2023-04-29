package mgit

import (
	"atilm/mgit/utilities"
	"os/exec"
	"path/filepath"
)

func Pull(baseDirectory string, args []string, printer utilities.ConsolePrinter) error {
	gitStatusItems, err := initializeStatusSlice(baseDirectory)
	if err != nil {
		return err
	}

	for i, statusItem := range gitStatusItems {
		if err := executeGitPull(baseDirectory, statusItem); err != nil {
			return err
		}

		newStatus, _ := executeStatusCommand(fullPath(baseDirectory, statusItem), statusItem)
		gitStatusItems[i] = newStatus
	}

	printStatusItems(gitStatusItems, printer)

	return nil
}

func fullPath(baseDirectory string, statusItem gitStatus) string {
	return filepath.Join(baseDirectory, statusItem.dirName)
}

func executeGitPull(baseDirectory string, statusItem gitStatus) error {
	pullCommand := exec.Command("git", "pull")
	pullCommand.Dir = fullPath(baseDirectory, statusItem)
	return pullCommand.Run()
}
