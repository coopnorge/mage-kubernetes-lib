package magekubernetes

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh" // sh contains helpful utility functions, like RunV
)

// Validate Kubernetes manifests for ArgCD applications related to this repository
func Validate() error {
	templates, err := renderTemplates()
	if err != nil {
		return err
	}
	mg.Deps(mg.F(kubeScore, templates))
	mg.Deps(mg.F(kubeConform, templates, "api-platform"))
	mg.Deps(Pallets)
	fmt.Println("Validation passed")
	return nil
}

// KubeScore runs kube-score on Kubernetes manifests for ArgCD applications related to this repository
func KubeScore() error {
	templates, err := renderTemplates()
	if err != nil {
		return err
	}
	return kubeScore(templates)
}

// KubeConform runs kubeconform on Kubernetes manifests for ArgCD applications related to this repository
func KubeConform() error {
	templates, err := renderTemplates()
	if err != nil {
		return err
	}
	return kubeConform(templates, "api-platform")
}

// Pallets validates the pallet files in the .pallet directory
func Pallets() error {
	pallets, err := listPalletFiles()
	if err != nil {
		return err
	}
	fmt.Println("Validating Pallets")
	return kubeConform(strings.Join(pallets, ","), "pallets")
}

// ArgoCDListApps show the apps related to this repository
func ArgoCDListApps() error {
	err := listArgoCDDeployments()
	if err != nil {
		return err
	}
	return nil
}

// ArgoCDDiff runs a diff between local changes and the current running state in ArgoCD
func ArgoCDDiff() error {
	repo, err := repoURL()
	if err != nil {
		return err
	}
	apps, err := getArgoCDDeployments(repo)
	if err != nil {
		return err
	}
	err = getArgoCDDiff(apps)
	if err != nil {
		return err
	}
	return nil
}

func kubeScore(paths string) error {
	if paths == "" {
		fmt.Printf("No files presented to kube-score, skipping")
		return nil
	}
	cmdOptions := []string{
		"score"}
	out, err := sh.Output("kube-score", append(cmdOptions, strings.Split(paths, ",")...)...)
	if err != nil {
		fmt.Printf("kube-score returned exit code: %d\n Output:\n %v Error:\n %v\n", sh.ExitStatus(err), out, err)
		return err
	}
	fmt.Println("kube-score passed")
	return nil
}

func kubeConform(paths string, schemaSelection string) error {
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
