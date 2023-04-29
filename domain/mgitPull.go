package mgit

import (
	"atilm/mgit/utilities"
	"os/exec"
)

func Pull(baseDirectory string, args []string, printer utilities.ConsolePrinter) error {
	gitStatusItems, err := CollectGitStatusFromSubdirectories(baseDirectory)
	if err != nil {
		return err
	}

	pullAndReport := func(status gitStatus) (gitStatus, error) {
		if err := executeGitPull(baseDirectory, status); err != nil {
			return status, err
		}

		return executeStatusCommand(fullPath(baseDirectory, status), status)
	}

	return parallelStatusUpdate(pullAndReport, gitStatusItems, printer)
}

func executeGitPull(baseDirectory string, statusItem gitStatus) error {
	pullCommand := exec.Command("git", "pull")
	pullCommand.Dir = fullPath(baseDirectory, statusItem)
	return pullCommand.Run()
}
