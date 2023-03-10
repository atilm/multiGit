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

func TestMain(m *testing.M) {
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
	// perhaps push a first commit in client 2 remotes before cloning into client 1
	logError(os.MkdirAll(testPath("client2"), os.ModeTemporary))
	runCommand(exec.Command("git", "clone", testPath("remote1"), testPath("client2/remote1")))
	runCommand(exec.Command("git", "clone", testPath("remote2"), testPath("client2/remote2")))

	// initial commit
	logError(createFile(testPath("client2/remote2/file1.txt")))
	runCommandInDir(exec.Command("git", "add", "."), testPath("client2/remote2"))
	runCommandInDir(exec.Command("git", "commit", "-m", "\"message 1\""), testPath("client2/remote2"))
	runCommandInDir(exec.Command("git", "push"), testPath("client2/remote2"))

	// create a reading client
	logError(os.MkdirAll(testPath("client1"), os.ModeTemporary))
	runCommand(exec.Command("git", "clone", testPath("remote1"), testPath("client1/remote1")))
	runCommand(exec.Command("git", "clone", testPath("remote2"), testPath("client1/remote2")))
	logError(os.MkdirAll(testPath("client1/nonGitDirectory"), os.ModeTemporary))

	// commit not in client1 copies
	logError(createFile(testPath("client2/remote2/file2.txt")))
	runCommandInDir(exec.Command("git", "add", "."), testPath("client2/remote2"))
	runCommandInDir(exec.Command("git", "commit", "-m", "\"message 1\""), testPath("client2/remote2"))
	runCommandInDir(exec.Command("git", "push"), testPath("client2/remote2"))

	result := m.Run()

	os.RemoveAll(testdataDirectory)

	os.Exit(result)
}

func whenTheStatusCommandIsExecuted(baseDirectory string) (string, error) {
	return mgit.ReportStatus(baseDirectory)
}

func TestStatusListsAllDirectoriesInGivenBaseDirectory(t *testing.T) {
	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))

	expectedResult := `01: remote1 (main) [ok]
02: remote2 (main) [1 to pull, 0 to push]
`

	if err != nil {
		t.Errorf("Expected no error. Got: %v", err)
	}

	if outString != expectedResult {
		t.Errorf("Actual: '%v' != Expected: '%v'", outString, expectedResult)
	}
}
