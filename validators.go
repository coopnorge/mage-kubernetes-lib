package magekubernetes

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/sh"
)

// kyvernoPoliciesValidator runs Kyverno validation on rendered Kubernetes manifests.
func kyvernoPoliciesValidator(paths string) error {
	if paths == "" {
		fmt.Println("No rendered templates provided for validation, skipping.")
		return nil
	}

	policyDir := "kyverno-policies"            // Directory where policies are stored
	templatePaths := strings.Split(paths, ",") // Split the input paths into a list

	// Prepare resource arguments for the Kyverno CLI
	resourceArgs := []string{}
	for _, templatePath := range templatePaths {
		templatePath = strings.TrimSpace(templatePath)
		if templatePath == "" {
			continue
		}
		resourceArgs = append(resourceArgs, "-r", templatePath)
	}

	if len(resourceArgs) == 0 {
		fmt.Println("No valid rendered templates found for validation.")
		return nil
	}

	cmdOptions := append([]string{"apply", policyDir, "-t", "--detailed-results", "--continue-on-fail"}, resourceArgs...)

	output, err := sh.Output("kyverno", cmdOptions...)
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("kyverno validation failed for policy %w", err)
	}

	fmt.Printf("Kyverno validation completed.\n")

	if strings.Contains(output, "violation") || strings.Contains(output, "failed") {
		return fmt.Errorf("kyverno validation issues found with policy: %s", output)
	}

	return nil
}

func kubeScoreValidator(paths string) error {
	if paths == "" {
		fmt.Printf("No files presented to kube-score, skipping")
		return nil
	}
	cmdOptions := []string{
		"score",
		"--ignore-container-cpu-limit", // disable requiring cpu limit
	}
	out, err := sh.Output("kube-score", append(cmdOptions, strings.Split(paths, ",")...)...)
	if err != nil {
		fmt.Printf("kube-score returned exit code: %d\n Output:\n %v Error:\n %v\n", sh.ExitStatus(err), out, err)
		return err
	}
	fmt.Println("kube-score passed")
	return nil
}

func kubeConformValidator(paths string, schemaSelection string) error {
	if paths == "" {
		fmt.Printf("No files presented to kubeconform, skipping")
		return nil
	}
	cmdOptions := []string{
		"-strict",
		"-verbose",
		"-schema-location", "default",
		"-schema-location", "https://raw.githubusercontent.com/coopnorge/kubernetes-schemas/main/" + schemaSelection + "/{{ .ResourceKind }}{{ .KindSuffix }}.json"}
	out, err := sh.Output("kubeconform", append(cmdOptions, strings.Split(paths, ",")...)...)
	if err != nil {
		fmt.Printf("kubeconform returned exit code: %d\n Output:\n %v Error:\n %v\n", sh.ExitStatus(err), out, err)
		return err
	}
	fmt.Println("kubeconform passed")
	return nil
}
