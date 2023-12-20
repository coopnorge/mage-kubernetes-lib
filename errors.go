package magekubernetes

import "fmt"

type (
	ErrUnableToParseRemoteURL struct {
		url string
	}
)

// Error implements error.
func (e *ErrUnableToParseRemoteURL) Error() string {
	return fmt.Sprintf("Unable to parse remote url %v", e.url)
}

func (e *ErrUnableToParseRemoteURL) Is(target error) bool {
	t, ok := target.(*ErrUnableToParseRemoteURL)
	if !ok {
		return false
	}
	if t.url != e.url {
		return false
	}
	return true
}

func NewErrUnableToParseRemoteURL(url string) *ErrUnableToParseRemoteURL {
	return &ErrUnableToParseRemoteURL{url: url}
}
