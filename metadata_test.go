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
		{"origin     http://github.com/coopnorge/helloworld.git (fetch)", "", NewErrUnableToParseRemoteURL("http://github.com/coopnorge/helloworld.git")},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.remote, tt.want)
		t.Run(testname, func(t *testing.T) {
			got, err := gitRemoteParser(tt.remote)
			if got != tt.want || !errors.Is(err, tt.err) {
				t.Errorf("\n got: %s,%v \nwant: %s,%v", got, err, tt.want, tt.err)
			}
		})
	}
}
