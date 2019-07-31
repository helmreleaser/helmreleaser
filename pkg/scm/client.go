package scm

import (
	"github.com/google/go-github/v27/github"
	"github.com/helmreleaser/helmreleaser/pkg/helmreleaser"
	"github.com/pkg/errors"
)

type SCMClient struct {
	GitHub *github.Client
}

func (s *SCMClient) PublishRelease(context helmreleaser.HelmReleaserContext, filename string) (string, error) {
	if s.GitHub != nil {
		return s.publishGitHubRelease(context, filename)
	}

	return "", errors.New("no scm provider found")
}
