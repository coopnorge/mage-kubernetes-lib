package magekubernetes

import (
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/sh"
	"gopkg.in/yaml.v3"
)

// ArgoCDAppHelm contains the info for rendering a helm file
type ArgoCDAppHelm struct {
	ReleaseName string   `yaml:"releaseName"`
	ValueFiles  []string `yaml:"valueFiles"`
}

// ArgoCDAppSource contains the info where to find the source for rendering
type ArgoCDAppSource struct {
	Helm    ArgoCDAppHelm `yaml:"helm"`
	Path    string        `yaml:"path"`
	RepoRUL string        `yaml:"repoURL"`
}

// ArgoCDAppSpec contains the app source
type ArgoCDAppSpec struct {
	Source  ArgoCDAppSource   `yaml:"source"`
	Sources []ArgoCDAppSource `yaml:"sources"`
}

// ArgoCDAppMetadata contains the app name
type ArgoCDAppMetadata struct {
	Name string `yaml:"name"`
}

// ArgoCDApp contains the spec and metadata of an app
type ArgoCDApp struct {
	Spec     ArgoCDAppSpec     `yaml:"spec"`
	Metadata ArgoCDAppMetadata `yaml:"metadata"`
}

func getArgoCDDeployments(repoURL string) ([]ArgoCDApp, error) {
	var argoCDAppList []ArgoCDApp
	env := map[string]string{}

	if token, ok := os.LookupEnv("ARGOCD_API_TOKEN"); ok {
		server, ok := env["ARGOCD_SERVER_NAME"]
		if !ok {
			return nil, fmt.Errorf("When using ARGOCD_API_TOKEN, you are also required to set ARGOCD_SERVER_NAME")
		}
		env["ARGOCD_API_TOKEN"] = token
		env["ARGOCD_SERVER_NAME"] = server
	} else {
		err := sh.Run("argocd", "context")
		if err != nil {
			fmt.Println("Make use $HOME/.argocd is correctly mounted or use ARGOCD_API_TOKEN env var")
			return nil, err
		}
	}
	// use label selector to quickly exclude pallet apps
	appYaml, err := sh.OutputWith(env, "argocd", "--grpc-web", "app", "list", "-r", repoURL, "-l", "component!=pallet-config", "-o", "yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(appYaml), &argoCDAppList)
	if err != nil {
		return nil, err
	}
	return argoCDAppList, nil
}

func getArgoCDDiff(apps []ArgoCDApp) error {
	env := map[string]string{"KUBECTL_EXTERNAL_DIFF": "dyff between --omit-header"}
	if token, ok := os.LookupEnv("ARGOCD_API_TOKEN"); ok {
		env["ARGOCD_API_TOKEN"] = token
	}
	for _, app := range apps {
		diff, err := sh.OutputWith(env, "argocd", "--loglevel", "error", "--grpc-web", "app", "diff", app.Metadata.Name, "--refresh", "--local", app.Spec.Source.Path)
		if sh.ExitStatus(err) == 2 {
			return err
		}
		fmt.Println("---- Diff of " + app.Metadata.Name + "  ----")
		fmt.Println(diff)
	}
	return nil
}

func listArgoCDDeployments() error {
	repo, err := repoURL()
	if err != nil {
		return err
	}
	apps, err := getArgoCDDeployments(repo)
	if err != nil {
		return err
	}
	for _, trackedDeployment := range apps {
		if trackedDeployment.Spec.Source.Helm.ReleaseName != "" {
			fmt.Println("---")
			fmt.Println("Found helm deployment with name: " + trackedDeployment.Metadata.Name)
			fmt.Println("  path: " + trackedDeployment.Spec.Source.Path)
			fmt.Println("  valueFiles: " + strings.Join(trackedDeployment.Spec.Source.Helm.ValueFiles, ", "))
		} else if _, err := os.Stat(trackedDeployment.Spec.Source.Path + "/kustomization.yaml"); err == nil {
			fmt.Println("---")
			fmt.Println("Found kustomize deployment with name: " + trackedDeployment.Metadata.Name)
			fmt.Println("  path: " + trackedDeployment.Spec.Source.Path)
		} else {
			fmt.Println("---")
			fmt.Println("Found plain deployment with name: " + trackedDeployment.Metadata.Name)
			fmt.Println("  path: " + trackedDeployment.Spec.Source.Path)
		}
	}
	return nil
}
