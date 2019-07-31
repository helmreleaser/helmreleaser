package helmreleaser

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func (h HelmReleaser) RenderString(context HelmReleaserContext, in string) (string, error) {
	t := template.New(in)
	tt := template.Must(t.Parse(in))
	output := new(bytes.Buffer)
	if err := tt.Execute(output, context); err != nil {
		return "", errors.Wrap(err, "failed to template string")
	}

	return output.String(), nil
}

func (h HelmReleaser) Render(context HelmReleaserContext, dir string) error {
	t := template.New("release")

	valuesYAML, err := ioutil.ReadFile(path.Join(dir, "values.yaml"))
	if err != nil {
		return errors.Wrap(err, "failed to read values.yaml")
	}
	valuesTemplate := template.Must(t.Parse(string(valuesYAML)))
	output := new(bytes.Buffer)
	if err := valuesTemplate.Execute(output, context); err != nil {
		return errors.Wrap(err, "failed to render values.yaml")
	}
	if err := ioutil.WriteFile(path.Join(dir, "values.yaml"), output.Bytes(), 0644); err != nil {
		return errors.Wrap(err, "failed to save updated values.yaml")
	}

	chartYAML, err := ioutil.ReadFile(path.Join(dir, "Chart.yaml"))
	if err != nil {
		return errors.Wrap(err, "failed to read Chart.yaml")
	}
	metadata := chart.Metadata{}
	if err := yaml.Unmarshal(chartYAML, &metadata); err != nil {
		return errors.Wrap(err, "failed to parse Chart.yaml")
	}
	metadata.Version = h.ChartVersion
	metadata.AppVersion = h.AppVersion
	chartYAML, err = yaml.Marshal(metadata)
	if err != nil {
		return errors.Wrap(err, "failed to update Chart.yaml")
	}
	chartTemplate := template.Must(t.Parse(string(chartYAML)))
	output = new(bytes.Buffer)
	if err := chartTemplate.Execute(output, context); err != nil {
		return errors.Wrap(err, "failed to render Chart.yaml")
	}
	if err = ioutil.WriteFile(path.Join(dir, "Chart.yaml"), output.Bytes(), 0644); err != nil {
		return errors.Wrap(err, "failed to save updated Chart.yaml")
	}

	return nil
}
