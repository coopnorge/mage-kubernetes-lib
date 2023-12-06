package magekubernetes

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Annotations contains the github repo slug
type Annotations struct {
	ProjectSlug string `yaml:"github.com/project-slug"`
}

// Metadata contains the annotations struct
type Metadata struct {
	Annotations Annotations `yaml:"annotations"`
}

// CatalogInfo contains the metadata of this repo
type CatalogInfo struct {
	Metadata Metadata `yaml:"metadata"`
}

func repoName() (string, error) {
	var catalogInfo CatalogInfo

	yamlFile, err := os.ReadFile("catalog-info.yaml")
	if err != nil {
		fmt.Printf("yamlFile.Get err #%v ", err)
		return "", err
	}
	err = yaml.Unmarshal(yamlFile, &catalogInfo)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
		return "", err
	}
	return catalogInfo.Metadata.Annotations.ProjectSlug, nil
}

func repoURL() (string, error) {
	repoName, err := repoName()
	if err != nil {
		return "", err
	}
	return "https://github.com/" + repoName, nil
}
