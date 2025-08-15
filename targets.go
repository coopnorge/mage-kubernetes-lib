package magekubernetes

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Validate Kubernetes manifests for ArgCD applications related to this repository
func Validate() error {
	skipKubeScore, err := getBoolEnv("SKIP_KUBE_SCORE", false)
	if err != nil {
		return err
	}

	if err := renderTemplatesAndValidate(!skipKubeScore, true, true); err != nil {
		return err
	}

	mg.Deps(Pallets)
	fmt.Println("Validation passed")
	return nil
}

// KubeScore runs kube-score on Kubernetes manifests for ArgCD applications related to this repository
func KubeScore() error {
	return renderTemplatesAndValidate(true, false, false)
}

// KubeConform runs kubeconform on Kubernetes manifests for ArgCD applications related to this repository
func KubeConform() error {
	return renderTemplatesAndValidate(false, false, true)
}

// Pallets validates the pallet files in the .pallet directory
func Pallets() error {
	pallets, err := listPalletFiles()
	if err != nil {
		return err
	}
	fmt.Println("Validating Pallets")
	return kubeConformValidator(strings.Join(pallets, ","), "pallets")
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

// ValidateKyverno runs render Kubernetes manifests and invokes validateKyvernoPolicies.
func ValidateKyverno() error {
	return renderTemplatesAndValidate(false, true, false)
}
