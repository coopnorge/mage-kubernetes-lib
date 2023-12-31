package magekubernetes

import (
	"strings"
	"testing"
)

func TestFailedKubeScore(t *testing.T) {
	paths := strings.Join([]string{"tests/templates/fail/templates/deployment.yaml", "tests/templates/fail/templates/service.yaml"}, ",")
	err := kubeScore(paths)
	if err == nil {
		t.Fatalf(`kubeScore(paths) should fail but passed`)
	}
}

func TestOKKubeScore(t *testing.T) {
	paths := strings.Join([]string{"tests/templates/ok/templates/configmap.yaml"}, ",")
	err := kubeScore(paths)
	if err != nil {
		t.Fatalf(`kubeScore(paths) should pass but failed with error %v`, err)
	}
}

func TestFailedKubeConform(t *testing.T) {
	paths := strings.Join([]string{"tests/templates/fail-schema/templates/deployment.yaml", "tests/templates/fail-schema/templates/service.yaml"}, ",")
	err := kubeConform(paths, "api-platform")
	if err == nil {
		t.Fatalf(`kubeConform(paths,"api-platform) should fail but passed`)
	}
}

func TestOKKubeConform(t *testing.T) {
	paths := strings.Join([]string{"tests/templates/ok/templates/configmap.yaml"}, ",")
	err := kubeConform(paths, "api-platform")
	if err != nil {
		t.Fatalf(`kubeConform(paths,"api-platform) should pass but failed with error %v`, err)
	}
}
