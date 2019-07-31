package chartserver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	chartserverv1beta1 "github.com/charthq/chartserver/pkg/apis/chartserver/v1beta1"
	"github.com/helmreleaser/helmreleaser/pkg/helmreleaser"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func CreateChartVersionSpec(h *helmreleaser.HelmReleaser, ctx *helmreleaser.HelmReleaserContext, namespace string, shasum string, downloadPath string) (*chartserverv1beta1.ChartVersion, error) {
	appVersion, err := h.RenderString(*ctx, h.AppVersion)
	if err != nil {
		return nil, errors.Wrap(err, "failed to template app version")
	}

	chartVersion, err := h.RenderString(*ctx, h.ChartVersion)
	if err != nil {
		return nil, errors.Wrap(err, "failed to template chart version")
	}

	description, err := h.RenderString(*ctx, h.Description)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render description")
	}

	spec := chartserverv1beta1.ChartVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%d-%d-%d", h.Name, ctx.Major, ctx.Minor, ctx.Patch),
			Namespace: namespace,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "chartserver.io/v1beta1",
			Kind:       "ChartVersion",
		},
		Spec: chartserverv1beta1.ChartVersionSpec{
			AppVersion:   appVersion,
			ChartVersion: chartVersion,
			Created:      time.Now().Format(time.RFC3339),
			Description:  description,
			Digest:       shasum,
			Home:         h.Home,
			Icon:         h.Icon,
			Maintainers:  h.Maintainers,
			Name:         h.Name,
			Sources:      h.Sources,
			URLs: []string{
				downloadPath,
			},
		},
	}

	return &spec, nil
}

func WriteSpecFile(spec *chartserverv1beta1.ChartVersion, filename string) error {
	output := new(bytes.Buffer)

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	if err := s.Encode(spec, output); err != nil {
		return errors.Wrap(err, "failed to marshal yaml")
	}

	if err := ioutil.WriteFile(filename, output.Bytes(), 0644); err != nil {
		return errors.Wrap(err, "failed to write yaml")
	}

	return nil
}
