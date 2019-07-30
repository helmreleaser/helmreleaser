package helmreleaser

import (
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type HelmReleaserContext struct {
	Major int64
	Minor int64
	Patch int64
}

// CreateContext will create a context that can be used to render
// The dir param should be the path in the git repo, not the temp directory
func CreateContext(dir string, scmToken string) (*HelmReleaserContext, error) {
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

	// scm := scm.NewGitHubClient(scmToken)
	// fmt.Printf("%#v\n", scm)

	helmReleaserContext := HelmReleaserContext{
		Major: latestTag.Major(),
		Minor: latestTag.Minor(),
		Patch: latestTag.Patch(),
	}

	return &helmReleaserContext, nil
}
