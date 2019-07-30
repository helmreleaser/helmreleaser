package helmreleaser

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type HelmReleaser struct {
	ChartVersion string  `yaml:"chartVersion"`
	AppVersion   string  `yaml:"appVersion"`
	Images       []Image `yaml:"images,omitempty"`

	Snapshot Snapshot `yaml:"snapshot,omitempty"`

	Archive *Archive `yaml:"archive,omitempty"`
}

type Archive struct {
	GitHub *GitHubArchive `yaml:"github,omitempty"`
}

type GitHubArchive struct {
	Owner string `yaml:"owner"`
	Name  string `yaml:"name"`
}

type Image struct {
	ID            string `yaml:"id"`
	ImageKey      string `yaml:"imageKey,omitempty"`
	TagKey        string `yaml:"tagKey,omitempty"`
	ImageTemplate string `yaml:"imageTemplate,omitempty"`
	TagTemplate   string `yaml:"tagTemplate,omitempty"`
}

type Snapshot struct {
	Images []Image `yaml:"images,omitempty"`
}

func CreateDefault() *HelmReleaser {
	return &HelmReleaser{
		ChartVersion: "{{ .Major }}.{{ .Minor }}.{{ .Patch }}",
		AppVersion:   "{{ .Major }}.{{ .Minor }}.{{ .Patch }}",
		Images:       []Image{},
		Snapshot:     Snapshot{},
		Archive: &Archive{
			GitHub: &GitHubArchive{
				Owner: "repo-owner-name",
				Name:  "repo-name",
			},
		},
	}
}

func (h *HelmReleaser) WriteToFile(filename string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(filename); err == nil {
			return errors.Errorf("file %s already exists", filename)
		}
	}
	b, err := yaml.Marshal(h)
	if err != nil {
		return errors.Wrap(err, "failed to marshal yaml")
	}

	if err := ioutil.WriteFile(filename, b, 0644); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}

	return nil
}

func ReadFromFile(filename string) (*HelmReleaser, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	helmReleaser := HelmReleaser{}
	if err := yaml.Unmarshal(b, &helmReleaser); err != nil {
		return nil, errors.Wrap(err, "failed to parse file")
	}

	return &helmReleaser, nil
}
