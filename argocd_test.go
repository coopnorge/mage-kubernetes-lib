package magekubernetes

import (
	"testing"
)

func TestGetArgoCDAuth(t *testing.T) {
	t.Setenv("ARGOCD_API_TOKEN", "token")
	t.Setenv("ARGOCD_SERVER", "server")
	options, err := getArgoCDAuth()
	want := []string{"--auth-token", "token", "--server", "server"}
	if !testStringSliceEq(options, want) || err != nil {
		t.Fatalf(`getArgoCDAuth() failed. \nWant: %s\n got: %s\n error: %v`, want, options, err)
	}
}
