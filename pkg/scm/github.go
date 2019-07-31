package scm

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/google/go-github/v27/github"
	"github.com/helmreleaser/helmreleaser/pkg/helmreleaser"
	"github.com/pkg/errors"
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

func (s *SCMClient) publishGitHubRelease(ctx helmreleaser.HelmReleaserContext, file string) error {
	_, assetName := path.Split(file)

	title := "HelmReleaser release"

	releaseData := &github.RepositoryRelease{
		Name:       github.String(title),
		TagName:    github.String(ctx.Tag),
		Body:       github.String("body"),
		Draft:      github.Bool(false),
		Prerelease: github.Bool(false),
	}

	release, _, err := s.GitHub.Repositories.GetReleaseByTag(context.Background(), ctx.GitRemote.Owner, ctx.GitRemote.Name, ctx.Tag)
	if err != nil {
		release, _, err = s.GitHub.Repositories.CreateRelease(context.Background(), ctx.GitRemote.Owner, ctx.GitRemote.Name, releaseData)
	} else {
		if release.GetBody() != "" {
			releaseData.Body = release.Body
		}

		release, _, err = s.GitHub.Repositories.EditRelease(context.Background(), ctx.GitRemote.Owner, ctx.GitRemote.Name, release.GetID(), releaseData)
	}

	f, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	releaseAsset, _, err := s.GitHub.Repositories.UploadReleaseAsset(context.Background(),
		ctx.GitRemote.Owner, ctx.GitRemote.Name, release.GetID(), &github.UploadOptions{
			Name: assetName,
		},
		f)
	if err != nil {
		return errors.Wrap(err, "failed to upload asset to github")
	}

	fmt.Printf("%s\n", *releaseAsset.BrowserDownloadURL)
	return nil
}
