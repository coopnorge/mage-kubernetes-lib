package magekubernetes

import (
	"os"
	"strings"
	"testing"
)

func TestFailedKubeScore(t *testing.T) {
	paths := strings.Join([]string{"tests/templates/fail/templates/deployment.yaml", "tests/templates/fail/templates/service.yaml"}, ",")
	err := kubeScoreValidator(paths)
	if err == nil {
		t.Fatalf(`kubeScore(paths) should fail but passed`)
	}
}

func TestOKKubeScore(t *testing.T) {
	paths := strings.Join([]string{"tests/templates/ok/templates/configmap.yaml"}, ",")
	err := kubeScoreValidator(paths)
	if err != nil {
		t.Fatalf(`kubeScore(paths) should pass but failed with error %v`, err)
	}
}

func TestFailedKubeConform(t *testing.T) {
	paths := strings.Join([]string{"tests/templates/fail-schema/templates/deployment.yaml", "tests/templates/fail-schema/templates/service.yaml"}, ",")
	err := kubeConformValidator(paths, "api-platform")
	if err == nil {
		t.Fatalf(`kubeConform(paths,"api-platform) should fail but passed`)
	}
}

func TestOKKubeConform(t *testing.T) {
	paths := strings.Join([]string{"tests/templates/ok/templates/configmap.yaml"}, ",")
	err := kubeConformValidator(paths, "api-platform")
	if err != nil {
		t.Fatalf(`kubeConform(paths,"api-platform) should pass but failed with error %v`, err)
	}
}

// Test for manifest files expected to fail Kyverno policy validation
func TestFailedValidateKyverno(t *testing.T) {
	path := "tests/templates/validate-fail/deployment-fail.yaml,tests/templates/validate-fail/deployment-fail.yaml,"
	err := kyvernoPoliciesValidator(path)
	if err == nil {
		t.Fatalf("Expected validation to fail for manifest %s, but it passed", path)
	}
}

// Test for manifest files expected to pass Kyverno policy validation
func TestOKValidateKyverno(t *testing.T) {
	path := "tests/templates/validate/deployment-ok.yaml,tests/templates/validate/deployment2-ok.yaml,"
	err := kyvernoPoliciesValidator(path)
	if err != nil {
		t.Fatalf("Expected validation to pass for manifest %s, but it failed with error: %v", path, err)
	}
}

func TestOKValidateEnvVars(t *testing.T) {
	appSource := &ArgoCDAppSource{
		Path: "tests/envvars/ok",
		Helm: ArgoCDAppHelm{
			ReleaseName: "test",
			ValueFiles:  []string{"values.yaml"},
		},
	}

	outDir, err := renderHelm(*appSource)
	if err != nil {
		t.Fatalf("renderHelm failed: %v", err)
	}
	if outDir == "" {
		t.Skip("renderHelm returned empty dir (likely missing chart path)")
	}

	files, err := listFilesInDirectory(outDir)
	if err != nil {
		t.Fatalf("listFilesInDirectory failed: %v", err)
	}

	if err := envVarsValidator(files); err != nil {
		t.Errorf("expected no duplicate env vars, got error: %v", err)
	}
}

func TestFailValidateEnvVars(t *testing.T) {
	appSource := &ArgoCDAppSource{
		Path: "tests/envvars/fail",
		Helm: ArgoCDAppHelm{
			ReleaseName: "test",
			ValueFiles:  []string{"values.yaml"},
		},
	}

	outDir, err := renderHelm(*appSource)
	if err != nil {
		t.Fatalf("renderHelm failed: %v", err)
	}
	if outDir == "" {
		t.Skip("renderHelm returned empty dir (likely missing chart path)")
	}

	files, err := listFilesInDirectory(outDir)
	if err != nil {
		t.Fatalf("listFilesInDirectory failed: %v", err)
	}

	if err := envVarsValidator(files); err == nil {
		t.Error("expected duplicate env vars to be detected, got nil")
	} else {
		t.Logf("duplicate env var detection succeeded: %v", err)
	}
}

func TestOKMultipleContextValidateEnvVars(t *testing.T) {
	files, err := listFilesInDirectory("tests/envvars/ok-multiple-env-contexts")
	if err != nil {
		t.Fatalf("listFilesInDirectory failed: %v", err)
	}
	if err := envVarsValidator(files); err != nil {
		t.Errorf("expected no duplicate env vars, got error: %v", err)
	}
}

func TestFailExistingEnvValidateEnvVars(t *testing.T) {
	appSource := &ArgoCDAppSource{
		Path: "tests/envvars/fail-existing-env",
		Helm: ArgoCDAppHelm{
			ReleaseName: "test",
			ValueFiles:  []string{"values.yaml"},
		},
	}

	outDir, err := renderHelm(*appSource)
	if err != nil {
		t.Fatalf("renderHelm failed: %v", err)
	}
	if outDir == "" {
		t.Skip("renderHelm returned empty dir (likely missing chart path)")
	}

	files, err := listFilesInDirectory(outDir)
	if err != nil {
		t.Fatalf("listFilesInDirectory failed: %v", err)
	}

	if err := envVarsValidator(files); err == nil {
		t.Error("expected duplicate env vars to be detected, got nil")
	} else {
		t.Logf("duplicate env var detection succeeded: %v", err)
	}
}

// Test for getBoolEnv returning expected values
func TestGetBoolEnv(t *testing.T) {
	os.Setenv("TEST_BOOL_1", "1")
	os.Setenv("TEST_BOOL_0", "0")
	os.Setenv("TEST_BOOL_TRUE", "true")
	os.Setenv("TEST_BOOL_FALSE", "false")
	os.Unsetenv("TEST_BOOL_UNSET")
	os.Setenv("TEST_BOOL_INVALID", "not_a_bool")

	tests := []struct {
		key      string
		def      bool
		expected bool
		wantErr  bool
	}{
		{"TEST_BOOL_1", false, true, false},
		{"TEST_BOOL_0", false, false, false},
		{"TEST_BOOL_TRUE", false, true, false},
		{"TEST_BOOL_FALSE", true, false, false},
		{"TEST_BOOL_UNSET", true, true, false},
		{"TEST_BOOL_UNSET", false, false, false},
		{"TEST_BOOL_INVALID", false, false, true},
	}

	for _, tt := range tests {
		got, err := getBoolEnv(tt.key, tt.def)
		if (err != nil) != tt.wantErr {
			t.Fatalf("Expected error=%v for key %s, got %v", tt.wantErr, tt.key, err)
		}
		if got != tt.expected {
			t.Fatalf("Expected %v for key %s, got %v", tt.expected, tt.key, got)
		}
	}
}
