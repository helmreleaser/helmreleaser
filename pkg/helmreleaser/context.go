package helmreleaser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

type HelmReleaserContext struct {
	Major int64
	Minor int64
	Patch int64

	Tag string

	GitRemote GitRemote

	changelog string
}

type GitRemote struct {
	Owner string
	Name  string
}

func (c HelmReleaserContext) GetChangelog() string {
	return c.changelog
}

// CreateContext will create a context that can be used to render
// The dir param should be the path in the git repo, not the temp directory
func (h HelmReleaser) CreateContext(dir string, scmToken string) (*HelmReleaserContext, error) {
	r, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open git repo")
	}

	tags, err := r.Tags()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list tags")
	}
	defer tags.Close()

	semverTags := []*semver.Version{}
	fullVersions := make(map[string]string)

	tags.ForEach(func(t *plumbing.Reference) error {
		tagNameSplit := strings.Split(t.Name().String(), "/")
		if len(tagNameSplit) < 3 {
			return errors.Errorf("unexpected tag format: %s", t.Name())
		}
		tagName := strings.Join(tagNameSplit[2:], "/")

		ver, err := semver.NewVersion(tagName)
		if err != nil {
			return errors.Wrap(err, "failed to parse semver tag")
		}

		fullVersions[ver.Original()] = t.Name().String()

		semverTags = append(semverTags, ver)
		return nil
	})

	if len(semverTags) == 0 {
		return nil, errors.New("no tags found in repo")
	}

	sort.Sort(semver.Collection(semverTags))
	for i := len(semverTags)/2 - 1; i >= 0; i-- {
		opp := len(semverTags) - 1 - i
		semverTags[i], semverTags[opp] = semverTags[opp], semverTags[i]
	}

	latestTag := semverTags[0]

	originRemote, err := r.Remote("origin")
	if err != nil {
		return nil, errors.New("no remote named 'origin' found")
	}
	endpoint, err := transport.NewEndpoint(originRemote.Config().URLs[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse git url")
	}
	splitPath := strings.Split(endpoint.Path, "/")

	helmReleaserContext := HelmReleaserContext{
		Tag:   latestTag.Original(),
		Major: latestTag.Major(),
		Minor: latestTag.Minor(),
		Patch: latestTag.Patch(),
		GitRemote: GitRemote{
			Owner: splitPath[0],
			Name:  splitPath[1],
		},
	}

	if h.Archive != nil {
		if h.Archive.GitHub != nil {
			if h.Archive.GitHub.Owner != "" {
				helmReleaserContext.GitRemote.Owner = h.Archive.GitHub.Owner
			}
			if h.Archive.GitHub.Name != "" {
				helmReleaserContext.GitRemote.Name = h.Archive.GitHub.Name
			}
		}
	}

	changelog, err := readChangelog(r, fullVersions[latestTag.Original()])
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a changelog")
	}
	helmReleaserContext.changelog = changelog

	return &helmReleaserContext, nil
}

func readChangelog(r *git.Repository, tag string) (string, error) {
	var hash plumbing.Hash

	tags, err := r.Tags()
	if err != nil {
		return "", errors.Wrap(err, "failed to parse git tags")
	}
	tags.ForEach(func(t *plumbing.Reference) error {
		if t.Name().String() == tag {
			hash = t.Hash()
		}
		return nil
	})

	logOptions := &git.LogOptions{
		From: hash,
	}
	logs, err := r.Log(logOptions)
	if err != nil {
		return "", errors.Wrap(err, "failed to read git log")
	}

	changelog := "## Changelog\n\n"
	logs.ForEach(func(commit *object.Commit) error {
		changelog = fmt.Sprintf("%s%s %s\n", changelog, commit.Hash.String()[0:7], commit.Message)
		return nil
	})

	return changelog, nil
}
