package scm

import (
	"context"

	"github.com/google/go-github/v27/github"
	"golang.org/x/oauth2"
)

func NewGitHubClient(token string) *SCMClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	scmClient := SCMClient{
		GitHub: client,
	}

	return &scmClient
}
