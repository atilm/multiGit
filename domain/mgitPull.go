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

	gitStatusItems = filterStatusItems(gitStatusItems, args)

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

		if number < 1 || number > statusItemCount {
			printer.PrintLines([]string{fmt.Sprintf("Repo index %s is not in range [1:%d].", arg, statusItemCount)})
			return ErrRepoIndex
		}
	}

	return nil
}

func filterStatusItems(items []gitStatus, args []string) []gitStatus {
	if len(args) == 0 {
		return items
	}

	argsSet := toSet(args)

	filteredList := make([]gitStatus, 0, len(items))
	for _, item := range items {
		if contains(argsSet, int(item.index+1)) {
			filteredList = append(filteredList, item)
		}
	}

	return filteredList
}

func toSet(args []string) map[int]struct{} {
	type void struct{}
	var member void

	set := make(map[int]struct{})
	for _, arg := range args {
		number, _ := strconv.Atoi(arg)
		set[number] = member
	}
	return set
}

func contains(set map[int]struct{}, element int) bool {
	_, exists := set[element]
	return exists
}

func executeGitPull(baseDirectory string, statusItem gitStatus) error {
	pullCommand := exec.Command("git", "pull")
	pullCommand.Dir = fullPath(baseDirectory, statusItem)
	return pullCommand.Run()
}
