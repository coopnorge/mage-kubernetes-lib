package magekubernetes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh" // sh contains helpful utility functions, like RunV
	"gopkg.in/yaml.v3"
)

func renderTemplate(app ArgoCDApp) (string, error) {
	if app.Spec.Source.Helm.ReleaseName != "" {
		return renderHelm(app.Spec.Source)
	} else if _, err := os.Stat(app.Spec.Source.Path + "/kustomization.yaml"); err == nil {
		return renderKustomize(app.Spec.Source.Path)
	}
	return app.Spec.Source.Path, nil
}

func renderHelm(source ArgoCDAppSource) (string, error) {
	dir, err := tempDir()
	if err != nil {
		return "", err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	err = os.Chdir(source.Path)
	if err != nil {
		return "", err
	}
	// temporary fix  until https://github.com/helm/helm/issues/7214 is fixed
	// again
	err = addHelmRepos("./")
	if err != nil {
		return "", err
	}
	fmt.Println("rendering helm templates to: " + dir)
	err = sh.Run("helm", "dependency", "build")
	if err != nil {
		return "", err
	}
	err = sh.Run("helm", "template",
		"-f", strings.Join(source.Helm.ValueFiles, ","),
		"--output-dir", dir,
		".")
	if err != nil {
		return "", err
	}
	err = os.Chdir(pwd)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func renderKustomize(path string) (string, error) {
	dir, err := tempDir()
	if err != nil {
		return "", err
	}
	fmt.Println("rendering kustomize templates: " + dir)
	err = sh.Run("kustomize", "build", path, "--output", dir)
	if err != nil {
		return "", err
	}
	return dir, nil
}

// Render templates to an temporary directory. Using a comma sep string here because
// mg. can only have int, str and bools as arguments
func renderTemplates() (string, error) {
	var files []string
	repo, err := repoURL()
	if err != nil {
		return "", err
	}
	apps, err := getArgoCDDeployments(repo)
	if err != nil {
		return "", err
	}
	for _, trackedDeployment := range apps {
		templates, err := renderTemplate(trackedDeployment)
		if err != nil {
			return "", err
		}
		tackedFiles, err := listFilesInDirectory(templates)
		if err != nil {
			return "", err
		}
		files = append(files, tackedFiles...)
	}
	return strings.Join(files, ","), nil
}

func addHelmRepos(path string) error {
	var chart HelmChart
	chartfile, err := os.ReadFile(filepath.Join(path, "Chart.yaml"))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(chartfile), &chart)
	if err != nil {
		return err
	}
	for _, dep := range chart.Dependencies {
		err := sh.Run("helm", "repo", "add", dep.Name, dep.Repository)
		if err != nil {
			return err
		}
	}
	return nil
}

// HelmChart contains all metadata of an helm chart
type HelmChart struct {
	Dependencies []HelmDependency `yaml:"dependencies"`
}

// HelmDependency contains a dependency of a helmchart
type HelmDependency struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	Repository string `yaml:"repository"`
	Alias      string `yaml:"alias"`
}
