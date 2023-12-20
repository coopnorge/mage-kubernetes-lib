package magekubernetes

import (
	"errors"
	"fmt"
	"testing"
)

func TestRepoURL(t *testing.T) {
	r, err := repoURL()
	want := "https://github.com/coopnorge/mage-kubernetes-lib"
	if r != want || err != nil {
		t.Fatalf(`repoURL() failed. \nWant: %s\n got: %s\n error: %v`, want, r, err)
	}
}

func TestGitRemoteParser(t *testing.T) {
	var tests = []struct {
		remote string
		want   string
		err    error
	}{
		{"origin     git@github.com:coopnorge/helloworld.git (fetch)", "https://github.com/coopnorge/helloworld", nil},
		{"origin     https://github.com/coopnorge/helloworld.git (fetch)", "https://github.com/coopnorge/helloworld", nil},
		{"origin     http://github.com/coopnorge/helloworld.git (fetch)", "", fmt.Errorf(unableToParseRemoteErr, "http://github.com/coopnorge/helloworld.git")},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.remote, tt.want)
		t.Run(testname, func(t *testing.T) {
			got, err := gitRemoteParser(tt.remote)
			if got == tt.want &&
				(errors.Is(err, tt.err) || // This is to compare if the error is of the same type, which
					// happen when both errors are nil,
					// The line below is to compare if the error message is the same as string
					err.Error() == tt.err.Error()) {
				return
			}
			t.Errorf("\n got: %s,%v \nwant: %s,%v", got, err, tt.want, tt.err)
		})
	}
}
