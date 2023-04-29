package mgit

import (
	"atilm/mgit/utilities"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

var (
	ErrNonNumericArg = errors.New("Not a numeric argument")
	ErrRepoIndex     = errors.New("Repo index out of bounds")
	ErrArgCount      = errors.New("Too many arguments")
)

func Pull(baseDirectory string, args []string, printer utilities.ConsolePrinter) error {
	gitStatusItems, err := CollectGitStatusFromSubdirectories(baseDirectory)
	if err != nil {
		return err
	}

	if err := validateArgs(args, len(gitStatusItems), printer); err != nil {
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

func validateArgs(args []string, statusItemCount int, printer utilities.ConsolePrinter) error {
	if len(args) == 0 {
		return nil
	}

	if len(args) > statusItemCount {
		printer.PrintLines([]string{"More arguments given than repos present."})
		return ErrArgCount
	}

	for _, arg := range args {
		number, err := strconv.Atoi(arg)

		if err != nil {
			printer.PrintLines([]string{fmt.Sprintf("Non-numeric argument %s found.", arg)})
			return ErrNonNumericArg
		}

		if number < 1 || number > len(args) {
			printer.PrintLines([]string{fmt.Sprintf("Repo index %s is not in range [1:%d].", arg, statusItemCount)})
			return ErrRepoIndex
		}
	}

	return nil
}

func executeGitPull(baseDirectory string, statusItem gitStatus) error {
	pullCommand := exec.Command("git", "pull")
	pullCommand.Dir = fullPath(baseDirectory, statusItem)
	return pullCommand.Run()
}
