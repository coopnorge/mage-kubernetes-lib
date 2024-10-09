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

	cmdOptions := []string{
		"--grpc-web",
		"app",
		"list",
		"-r", repoURL,
		"-l", "component!=pallet-config", // use label selector to quickly exclude pallet apps
		"-o", "yaml",
	}

	authOptions, err := getArgoCDAuth()
	if err != nil {
		return nil, err
	}
	appYaml, err := sh.Output("argocd", append(cmdOptions, authOptions...)...)
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
	authOptions, err := getArgoCDAuth()
	if err != nil {
		return err
	}
	for _, app := range apps {
		cmdOptions := []string{
			"--loglevel", "error",
			"--grpc-web",
			"app",
			"diff",
			app.Metadata.Name,
			"--refresh",
			"--local", app.Spec.Source.Path,
			"--server-side-generate",
		}
		diff, err := sh.OutputWith(env, "argocd", append(cmdOptions, authOptions...)...)

		// `argocd diff` returns the following exit codes:
		// 		- 2 on general errors,
		//		- 1 when a diff is found, and
		//		- 0 when no diff
		// In addition, it also returns additional exit codes 11, 12, 13 and 20
		errCode := sh.ExitStatus(err)
		if errCode != 0 && errCode != 1 {
			return err
		}
		fmt.Println("---- Diff of " + app.Metadata.Name + "  ----")
		fmt.Println(diff)
	}
	return nil
}

func getArgoCDAuth() ([]string, error) {
	authOptions := []string{}
	if token, ok := os.LookupEnv("ARGOCD_API_TOKEN"); ok {
		server, ok := os.LookupEnv("ARGOCD_SERVER")
		if !ok {
			return nil, fmt.Errorf("When using ARGOCD_API_TOKEN, you are also required to set ARGOCD_SERVER")
		}
		authOptions = append(authOptions, "--auth-token", token, "--server", server)
	} else {
		err := sh.Run("argocd", "context")
		if err != nil {
			fmt.Println("Make use '$HOME/.argocd' is mounted to /root/.config/ or use ARGOCD_API_TOKEN and ARGOCD_SERVER environment variable")
			return nil, err
		}
	}
	return authOptions, nil
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
		} else if isKustomizeDir(trackedDeployment.Spec.Source.Path) {
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
