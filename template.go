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
	fmt.Printf("Preparing to render template for app: %s\n", app.Metadata.Name)
	if app.Spec.Source.Helm.ReleaseName != "" {
		fmt.Printf("Rendering helm release %s\n", app.Spec.Source.Helm.ReleaseName)
		return renderHelm(app.Spec.Source)
	} else if isKustomizeDir(app.Spec.Source.Path) {
		fmt.Printf("Rendering kustomize %s\n", app.Spec.Source.Path)
		return renderKustomize(app.Spec.Source.Path)
	}
	fmt.Printf("Rendering template %s\n", app.Spec.Source.Path)
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
		// dont error here, it seems we cannot
		// find the directory so we dont render templates
		// cause could be wrong configuration of the argocd app
		// or the new config is not yet on the main branch
		fmt.Printf("Directory %s not found. Skipping rendering manifests.\n", source.Path)
		return "", nil
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
		"--skip-tests",
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
	fmt.Println("rendering templates for repo: " + repo)
	if err != nil {
		return "", err
	}
	apps, err := getArgoCDDeployments(repo)
	if err != nil {
		return "", fmt.Errorf("getting ArgoCD deployments failed: %w", err)
	}
	for _, trackedDeployment := range apps {
		templates, err := renderTemplate(trackedDeployment)
		if err != nil {
			return "", fmt.Errorf("rendering templates failed for %s: %w", trackedDeployment, err)
		}
		fmt.Println("listing files in templates directory: " + templates)
		tackedFiles, err := listFilesInDirectory(templates)
		if err != nil {
			return "", fmt.Errorf("listing files failed for %s: %w", trackedDeployment, err)
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
		if strings.HasPrefix(dep.Repository, "oci://") {
			fmt.Println("skipping repo add for oci repository: " + dep.Repository)
			continue
		}

		fmt.Println("adding repo: " + dep.Repository)

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
