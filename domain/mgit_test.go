package mgit_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func logError(err error) {
	if err != nil {
		fmt.Print(err)
	}
}

func runCommand(command *exec.Cmd) {
	err := command.Run()
	logError(err)
}

func TestMain(m *testing.M) {
	logError(os.MkdirAll("testdata/remote1", os.ModeTemporary))
	logError(os.MkdirAll("testdata/remote2", os.ModeTemporary))
	logError(os.MkdirAll("testdata/client1", os.ModeTemporary))
	logError(os.MkdirAll("testdata/client2", os.ModeTemporary))
	runCommand(exec.Command("git", "--bare", "init", "testdata/remote1"))
	runCommand(exec.Command("git", "--bare", "init", "testdata/remote2"))
	runCommand(exec.Command("git", "--bare", "init", "testdata/remote2"))
	runCommand(exec.Command("git", "--bare", "init", "testdata/remote2"))
	runCommand(exec.Command("git", "clone", "testdata/remote1", "testdata/client1/remote1"))
	runCommand(exec.Command("git", "clone", "testdata/remote2", "testdata/client1/remote2"))
	runCommand(exec.Command("git", "clone", "testdata/remote1", "testdata/client2/remote1"))
	runCommand(exec.Command("git", "clone", "testdata/remote2", "testdata/client2/remote2"))

	result := m.Run()

	os.RemoveAll("testdata/")

	os.Exit(result)
}

func TestADummy(t *testing.T) {
	t.Fail()
}
