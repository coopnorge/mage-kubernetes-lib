package magekubernetes

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
)

const (
	unableToParseRemoteErr = "unable to parse remote url: %s"
)

func repoURL() (string, error) {
	repo, err := git.PlainOpen("./")
	if err != nil {
		return "", err
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		return "", err
	}
	return gitRemoteParser(remote.String())
}

func gitRemoteParser(remote string) (string, error) {
	url := strings.Fields(remote)[1]
	if strings.HasPrefix(url, "https://") {
		return strings.TrimSuffix(url, ".git"), nil
	} else if strings.HasPrefix(url, "git@") {
		toHTTPS := "https://github.com/" + strings.Split(url, ":")[1]
		return strings.TrimSuffix(toHTTPS, ".git"), nil
	}
	return "", fmt.Errorf(unableToParseRemoteErr, url)
}
