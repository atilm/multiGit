package mgit_test

import (
	mgit "atilm/mgit/domain"
	"testing"
)

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
	outString, err = whenThePullCommandIsExecuted(testPath("client1"))
	thenThereIsNoError(err, t)
	expectedResult = `01: remote1 (main) [ok]
02: remote2 (main) [ok]`
	thenTheOutputIs(expectedResult, outString, t)
}

func TestPullWithInvalidArgumentsReturnsAnError(t *testing.T) {
	teardown := givenAnEnvironmentWithTwoClientsAndTwoRemotes(t)
	defer teardown(t)

	testCases := map[string]struct {
		arguments      []string
		expectedError  error
		expectedOutput string
	}{
		"non-numeric":          {arguments: []string{"n"}, expectedError: mgit.ErrNonNumericArg, expectedOutput: "Non-numeric argument n found."},
		"out of uppper bounds": {arguments: []string{"3"}, expectedError: mgit.ErrRepoIndex, expectedOutput: "Repo index 3 is not in range [1:2]."},
		"out of lower bounds":  {arguments: []string{"0"}, expectedError: mgit.ErrRepoIndex, expectedOutput: "Repo index 0 is not in range [1:2]."},
		"second arg invalid":   {arguments: []string{"1", "3"}, expectedError: mgit.ErrRepoIndex, expectedOutput: "Repo index 3 is not in range [1:2]."},
		"too many arguments":   {arguments: []string{"1", "2", "3"}, expectedError: mgit.ErrArgCount, expectedOutput: "More arguments given than repos present."},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			outString, err := whenThePullCommandIsExecuted(testPath("client1"), testCase.arguments...)
			thenThereIsAnError(err, testCase.expectedError, t)
			thenTheOutputIs(testCase.expectedOutput, outString, t)
		})
	}
}
