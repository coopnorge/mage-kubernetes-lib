package magekubernetes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func renderTemplate(app ArgoCDApp) (string, error) {
	defer since("renderTemplate", time.Now(), map[string]any{
		"app": app.Metadata.Name,
	})

	infof("Preparing to render template for app name=%q sourcePath=%q sourceType=%q releaseName=%q",
		app.Metadata.Name, app.Spec.Source.Path, app.Status.SourceType, app.Spec.Source.Helm.ReleaseName)

	if app.Spec.Source.Helm.ReleaseName != "" || app.Status.SourceType == "Helm" {
		infof("Selected renderer=helm path=%q", app.Spec.Source.Path)
		return renderHelm(app.Spec.Source)
	} else if isKustomizeDir(app.Spec.Source.Path) {
		infof("Selected renderer=kustomize path=%q", app.Spec.Source.Path)
		return renderKustomize(app.Spec.Source.Path)
	}
	infof("Selected renderer=raw path=%q", app.Spec.Source.Path)
	return app.Spec.Source.Path, nil
}

func renderHelm(source ArgoCDAppSource) (string, error) {
	defer since("renderHelm", time.Now(), map[string]any{
		"path": source.Path,
	})
	dir, err := tempDir()
	if err != nil {
		errorf("tempDir failed: %v", err)
		return "", err
	}
	infof("Helm render tempDir=%q", dir)

	pwd, err := os.Getwd()
	if err != nil {
		errorf("Getwd failed: %v", err)
		return "", err
	}
	debugf("cwd before chdir=%q", pwd)

	if err := os.Chdir(source.Path); err != nil {
		// dont error here, it seems we cannot
		// find the directory so we dont render templates
		// cause could be wrong configuration of the argocd app
		// or the new config is not yet on the main branch
		warnf("Directory %q not found; skipping helm rendering (likely misconfig or not yet merged).", source.Path)
		return "", nil
	}

	// restore original working directory when done
	defer func() {
		if chErr := os.Chdir(pwd); chErr != nil {
			errorf("failed to restore working directory to %q: %v", pwd, chErr)
		} else {
			debugf("restored working directory to %q", pwd)
		}
	}()

	// temporary fix until https://github.com/helm/helm/issues/7214 is fixed
	// again
	if err := addHelmRepos("./"); err != nil {
		errorf("addHelmRepos failed: %v", err)
		return "", err
	}

	infof("Rendering helm templates to %q", dir)

	if err := runLogged("helm", "dependency", "build"); err != nil {
		return "", err
	}

	values := strings.Join(source.Helm.ValueFiles, ",")
	if values != "" {
		debugf("Helm value files: %q", values)
	}

	args := []string{"template", "--skip-tests", "--output-dir", dir, "."}

	for _, vf := range source.Helm.ValueFiles {
		args = append(args, "-f", vf)
	}

	if strings.TrimSpace(source.Helm.Values) != "" {
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("inline-values-%d.yaml", time.Now().UnixNano()))
		debugf("Writing helm values to a temp file %q: %q", tmp, source.Helm.Values)
		if err := os.WriteFile(tmp, []byte(source.Helm.Values), 0o600); err != nil {
			return "", fmt.Errorf("write inline values: %w", err)
		}
		args = append(args, "-f", tmp)
	}

	if err := runLogged("helm", args...); err != nil {
		return "", err
	}

	infof("Helm rendering complete outputDir=%q", dir)
	return dir, nil
}

func renderKustomize(path string) (string, error) {
	defer since("renderKustomize", time.Now(), map[string]any{
		"path": path,
	})
	dir, err := tempDir()
	if err != nil {
		errorf("tempDir failed: %v", err)
		return "", err
	}
	infof("Rendering kustomize templates outputDir=%q", dir)

	if err := runLogged("kustomize", "build", path, "--output", dir); err != nil {
		return "", err
	}
	infof("Kustomize rendering complete outputDir=%q", dir)
	return dir, nil
}

// Render templates to a temporary directory and validates using the selected
// validators
func renderTemplatesAndValidate(validateKubeScore bool, validateKyverno bool, validateKubeConform bool) error {
	defer since("renderTemplatesAndValidate", time.Now(), nil)

	repo, err := repoURL()
	infof("Rendering templates and validating repo=%q", repo)
	if err != nil {
		errorf("repoURL failed: %v", err)
		return err
	}

	apps, err := getArgoCDDeployments(repo)
	if err != nil {
		errorf("getArgoCDDeployments failed: %v", err)
		return fmt.Errorf("getting ArgoCD deployments failed: %w", err)
	}
	infof("Found %d ArgoCD apps to process", len(apps))

	for i, trackedDeployment := range apps {
		infof("(%d/%d) Start app name=%q path=%q sourceType=%q",
			i+1, len(apps),
			trackedDeployment.Metadata.Name,
			trackedDeployment.Spec.Source.Path,
			trackedDeployment.Status.SourceType,
		)

		templates, err := renderTemplate(trackedDeployment)
		if err != nil {
			errorf("renderTemplate failed for app=%q: %v", trackedDeployment.Metadata.Name, err)
			return fmt.Errorf("rendering templates failed for %v: %w", trackedDeployment, err)
		}

		if templates == "" {
			warnf("Templates path is empty for app=%q; skipping file listing", trackedDeployment.Metadata.Name)
			continue
		}

		infof("Listing files in templates directory dir=%q app=%q", templates, trackedDeployment.Metadata.Name)
		tackedFiles, err := listFilesInDirectory(templates)
		if err != nil {
			errorf("listFilesInDirectory failed dir=%q app=%q: %v", templates, trackedDeployment.Metadata.Name, err)
			return fmt.Errorf("listing files failed for %v: %w", trackedDeployment, err)
		}
		debugf("Discovered %d files in dir=%q", len(tackedFiles), templates)

		if validateKubeScore {
			debugf("Running kubeScoreValidator on %d files of %q", len(tackedFiles), trackedDeployment.Metadata.Name)
			if err := kubeScoreValidator(strings.Join(tackedFiles, ",")); err != nil {
				return err
			}
		}

		if validateKyverno {
			debugf("Running kyvernoPoliciesValidator on %d files of %q", len(tackedFiles), trackedDeployment.Metadata.Name)
			if err := kyvernoPoliciesValidator(strings.Join(tackedFiles, ",")); err != nil {
				return err
			}
		}

		if validateKubeConform {
			debugf("Running kubeConformValidator on %d files of %q", len(tackedFiles), trackedDeployment.Metadata.Name)
			if err := kubeConformValidator(strings.Join(tackedFiles, ","), "api-platform"); err != nil {
				return err
			}
		}
	}

	return nil
}

func addHelmRepos(path string) error {
	defer since("addHelmRepos", time.Now(), map[string]any{"path": path})

	var chart HelmChart
	chartfile, err := os.ReadFile(filepath.Join(path, "Chart.yaml"))
	if err != nil {
		errorf("reading Chart.yaml failed path=%q: %v", path, err)
		return err
	}
	debugf("Read Chart.yaml bytes=%d", len(chartfile))

	if err := yaml.Unmarshal(chartfile, &chart); err != nil {
		errorf("yaml unmarshal Chart.yaml failed: %v", err)
		return err
	}
	infof("Chart dependencies found: %d", len(chart.Dependencies))

	for _, dep := range chart.Dependencies {
		if strings.HasPrefix(dep.Repository, "oci://") {
			infof("Skipping helm repo add for OCI repository name=%q repo=%q", dep.Name, dep.Repository)
			continue
		}

		infof("Adding helm repo name=%q repo=%q", dep.Name, dep.Repository)
		if err := runLogged("helm", "repo", "add", dep.Name, dep.Repository); err != nil {
			return err
		}
	}
	infof("Helm repos for path %q processed successfully", path)
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
