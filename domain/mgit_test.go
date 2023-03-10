package mgit_test

import (
	mgit "atilm/mgit/domain"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func logError(err error) {
	if err != nil {
		fmt.Print(err)
	}
}

func runCommand(command *exec.Cmd) {
	output, err := command.CombinedOutput()

	if err != nil {
		fmt.Print(string(output))
	}

	logError(err)
}

func runCommandInDir(command *exec.Cmd, directory string) {
	command.Dir = directory
	runCommand(command)
}

var testdataDirectory string

func testPath(path string) string {
	return filepath.Join(testdataDirectory, path)
}

func createFile(filePath string) error {
	file, err := os.Create(filePath)
	file.Close()
	return err
}

func givenAnEnvironmentWithTwoClientsAndTworRemotes(t *testing.T) func(t *testing.T) {
	var err error
	testdataDirectory, err = ioutil.TempDir("", "mgit")
	if err != nil {
		os.Exit(1)
	}

	fmt.Println(testdataDirectory)

	// init bare "remote" repositories
	logError(os.MkdirAll(testPath("remote1"), os.ModeTemporary))
	logError(os.MkdirAll(testPath("remote2"), os.ModeTemporary))
	runCommand(exec.Command("git", "--bare", "init", testPath("remote1")))
	runCommand(exec.Command("git", "--bare", "init", testPath("remote2")))

	// create a writing client
	logError(os.MkdirAll(testPath("client2"), os.ModeTemporary))
	runCommand(exec.Command("git", "clone", testPath("remote1"), testPath("client2/remote1")))
	runCommand(exec.Command("git", "clone", testPath("remote2"), testPath("client2/remote2")))

	// initial commits
	addFileCommitAndPush(testPath("client2/remote1"), "file1.txt")
	addFileCommitAndPush(testPath("client2/remote2"), "file1.txt")

	// create a reading client
	logError(os.MkdirAll(testPath("client1"), os.ModeTemporary))
	runCommand(exec.Command("git", "clone", testPath("remote1"), testPath("client1/remote1")))
	runCommand(exec.Command("git", "clone", testPath("remote2"), testPath("client1/remote2")))
	logError(os.MkdirAll(testPath("client1/nonGitDirectory"), os.ModeTemporary))

	return func(t *testing.T) {
		os.RemoveAll(testdataDirectory)
	}
}

func addFileCommitAndPush(repositoryPath string, fileName string) {
	filePath := filepath.Join(repositoryPath, fileName)
	logError(createFile(filePath))
	runCommandInDir(exec.Command("git", "add", "."), repositoryPath)
	runCommandInDir(exec.Command("git", "commit", "-m", "\"message\""), repositoryPath)
	runCommandInDir(exec.Command("git", "push"), repositoryPath)
}

func whenTheStatusCommandIsExecuted(baseDirectory string) (string, error) {
	return mgit.ReportStatus(baseDirectory)
}

func whenAFileIsAddedTo(directory string, fileName string) {
	filePath := filepath.Join(directory, fileName)
	createFile(filePath)
}

func thenThereIsNoError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("Expected no error. Got: %v", err)
	}
}

func thenTheOutputIs(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Actual: '%v' != Expected: '%v'", actual, expected)
	}
}

func TestStatusListsAllDirectoriesInGivenBaseDirectory(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTworRemotes(t)
	defer teardown(t)

	// when someone else pushes something to remote2
	addFileCommitAndPush(testPath("client2/remote2"), "file2.txt")
	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))

	thenThereIsNoError(err, t)

	expectedResult := `01: remote1 (main) [ok]
02: remote2 (main) [1 to pull, 0 to push]
`
	thenTheOutputIs(expectedResult, outString, t)
}

func TestStatusReportsUntrackedFiles(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTworRemotes(t)
	defer teardown(t)

	whenAFileIsAddedTo(testPath("client1/remote1"), "newFile.txt")
	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))

	thenThereIsNoError(err, t)

	expectedOutput := `01: remote1 (main) * [ok]
02: remote2 (main) [ok]
`
	thenTheOutputIs(expectedOutput, outString, t)
}
