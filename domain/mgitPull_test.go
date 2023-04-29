package mgit_test

import "testing"

func TestPullWithoutArgumentsPullsAllRepos(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTwoRemotes(t)
	defer teardown(t)

	// when someone else pushes to both remotes
	addFileCommitAndPush(testPath("client2/remote1"), "file2.txt")
	addFileCommitAndPush(testPath("client2/remote2"), "file2.txt")

	// then status shows that there are changes to pull
	outString, err := whenTheStatusCommandIsExecuted(testPath("client1"))
	thenThereIsNoError(err, t)
	expectedResult := `01: remote1 (main) [1 to pull, 0 to push]
02: remote2 (main) [1 to pull, 0 to push]`
	thenTheOutputIs(expectedResult, outString, t)

	// mgit pull updates all repos
	outString, err = whenThePullCommandIsExecutedWithoutArgs(testPath("client1"))
	thenThereIsNoError(err, t)
	expectedResult = `01: remote1 (main) [ok]
02: remote2 (main) [ok]`
	thenTheOutputIs(expectedResult, outString, t)
}
