package chartserver

import (
	"fmt"
	"time"

	chartserverv1beta1 "github.com/charthq/chartserver/pkg/apis/chartserver/v1beta1"
	"github.com/helmreleaser/helmreleaser/pkg/helmreleaser"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PublishChartVersion(h *helmreleaser.HelmReleaser, ctx *helmreleaser.HelmReleaserContext, namespace string, shasum string, downloadPath string, outputFile string) (string, error) {
	appVersion, err := h.RenderString(*ctx, h.AppVersion)
	if err != nil {
		return "", errors.Wrap(err, "failed to template app version")
	}

	chartVersion, err := h.RenderString(*ctx, h.ChartVersion)
	if err != nil {
		return "", errors.Wrap(err, "failed to template chart version")
	}

	description, err := h.RenderString(*ctx, h.Description)
	if err != nil {
		return "", errors.Wrap(err, "failed to render description")
	}

	spec := chartserverv1beta1.ChartVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%d-%d-%d", h.Name, ctx.Major, ctx.Minor, ctx.Patch),
			Namespace: namespace,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "chart.sh/v1beta1",
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
			URLs:         h.URLs,
		},
	}

	b, err := yaml.Marshal(spec)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal chartversion spec")
	}

	return string(b), nil
}
