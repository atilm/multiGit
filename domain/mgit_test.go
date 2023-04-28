package mgit_test

import (
	"testing"
)

func TestStatusListsCommitsToPull(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTwoRemotes(t)
	defer teardown(t)

	// when someone else pushes something to remote2
	addFileCommitAndPush(testPath("client2/remote2"), "file2.txt")
	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))

	thenThereIsNoError(err, t)

	expectedResult := `01: remote1 (main) [ok]
02: remote2 (main) [1 to pull, 0 to push]`
	thenTheOutputIs(expectedResult, outString, t)
}

func TestStatusListsCommitsToPush(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTwoRemotes(t)
	defer teardown(t)

	whenAFileIsEdited(testPath("client1/remote2/file1.txt"), "a new line\n")
	whenAllChangesAreStaged(testPath("client1/remote2"))
	whenAllChangesAreCommitted(testPath("client1/remote2"))
	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))

	thenThereIsNoError(err, t)

	expectedResult := `01: remote1 (main) [ok]
02: remote2 (main) [0 to pull, 1 to push]`
	thenTheOutputIs(expectedResult, outString, t)
}

func TestStatusReportsUntrackedFiles(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTwoRemotes(t)
	defer teardown(t)

	whenAFileIsAddedTo(testPath("client1/remote1"), "newFile.txt")
	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))

	thenThereIsNoError(err, t)

	expectedOutput := `01: remote1 (main) * [ok]
02: remote2 (main) [ok]`
	thenTheOutputIs(expectedOutput, outString, t)
}

func TestStatusReportsUnstagedChanges(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTwoRemotes(t)
	defer teardown(t)

	whenAFileIsEdited(testPath("client1/remote2/file1.txt"), "text to append")
	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))

	thenThereIsNoError(err, t)

	expectedOutput := `01: remote1 (main) [ok]
02: remote2 (main) * [ok]`
	thenTheOutputIs(expectedOutput, outString, t)
}

func TestStatusReportsUncommmittedStagedChanges(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTwoRemotes(t)
	defer teardown(t)

	whenAFileIsEdited(testPath("client1/remote2/file1.txt"), "text to append")
	whenAllChangesAreStaged(testPath("client1/remote2/"))
	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))

	thenThereIsNoError(err, t)

	expectedOutput := `01: remote1 (main) [ok]
02: remote2 (main) * [ok]`
	thenTheOutputIs(expectedOutput, outString, t)
}

func TestStatusListsCommitsToPushAndPullWithLocalchanges(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTwoRemotes(t)
	defer teardown(t)

	// when someone else pushes something to remote2
	addFileCommitAndPush(testPath("client2/remote2"), "file2.txt")

	// when there are unpushed commits
	whenAFileIsEdited(testPath("client1/remote2/file1.txt"), "a new line\n")
	whenAllChangesAreStaged(testPath("client1/remote2"))
	whenAllChangesAreCommitted(testPath("client1/remote2"))

	// when there are uncommitted changes
	whenAFileIsAddedTo(testPath("client1/remote2"), "aNewFile.txt")

	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))

	thenThereIsNoError(err, t)

	expectedResult := `01: remote1 (main) [ok]
02: remote2 (main) * [1 to pull, 1 to push]`
	thenTheOutputIs(expectedResult, outString, t)
}
