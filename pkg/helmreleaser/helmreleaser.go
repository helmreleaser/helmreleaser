package helmreleaser

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type HelmReleaser struct {
	ChartVersion string   `yaml:"chartVersion,omitempty"`
	AppVersion   string   `yaml:"appVersion,omitempty"`
	Description  string   `yaml:"description,omitempty"`
	Home         string   `yaml:"home,omitempty"`
	Icon         string   `yaml:"icon,omitempty"`
	Maintainers  []string `yaml:"maintainers,omitempty"`
	Name         string   `yaml:"name,omitempty"`
	Sources      []string `yaml:"sources,omitempty"`
	URLs         []string `yaml:"urls,omitempty"`

	Images []Image `yaml:"images,omitempty"`

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

func (h *HelmReleaser) MergeValuesFromChart(filename string) error {
	chartYAML, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "failed to read to merge Chart.yaml")
	}

	metadata := chart.Metadata{}
	if err := yaml.Unmarshal(chartYAML, &metadata); err != nil {
		return errors.Wrap(err, "failed to parse to merge Chart.yaml")
	}

	h.Sources = metadata.GetSources()
	if h.Name == "" {
		h.Name = metadata.GetName()
	}
	h.Maintainers = []string{}
	for _, maintainer := range metadata.GetMaintainers() {
		h.Maintainers = append(h.Maintainers, maintainer.String())
	}
	h.Icon = metadata.GetIcon()
	h.Home = metadata.GetHome()
	if h.Description == "" {
		h.Description = metadata.GetDescription()
	}
	if h.AppVersion == "" {
		h.AppVersion = metadata.GetAppVersion()
	}
	if h.ChartVersion == "" {
		h.ChartVersion = metadata.GetVersion()
	}

	return nil
}
