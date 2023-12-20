package magekubernetes

import (
	"strings"

	"github.com/go-git/go-git/v5"
)

func repoURL() (string, error) {
	repo, err := git.PlainOpen("./")
	if err != nil {
		return "", err
	}
	repoURL, err := repo.Remote("origin")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(repoURL.String(), ".git"), nil
}
