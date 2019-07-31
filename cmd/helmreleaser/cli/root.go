package cli

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/helmreleaser/helmreleaser/pkg/chartserver"
	"github.com/helmreleaser/helmreleaser/pkg/helmreleaser"
	"github.com/helmreleaser/helmreleaser/pkg/logger"
	"github.com/helmreleaser/helmreleaser/pkg/scm"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/helm/pkg/chartutil"
)

func RootCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "helmreleaser",
		Short: "",
		Long:  `.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			logger := logger.NewLogger()

			logger.Info("")

			// read the config
			p := path.Join(v.GetString("chart-dir"), ".helmreleaser.yaml")
			logger.Info("Reading config from %s", p)
			helmReleaser, err := helmreleaser.ReadFromFile(p)
			if err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}

			cp := path.Join(v.GetString("chart-dir"), "Chart.yaml")
			logger.Info("Merging values from Chart.yaml")
			if err := helmReleaser.MergeValuesFromChart(cp); err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}

			// copy to a temp directory
			dir, err := ioutil.TempDir("", "helmreleaser")
			if err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}
			defer func() {
				os.RemoveAll(dir)
			}()

			logger.Info("Creating workspace directory for chart")
			if err := copy.Copy(v.GetString("chart-dir"), dir); err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}

			// render
			context, err := helmReleaser.CreateContext(v.GetString("chart-dir"), os.Getenv("GITHUB_TOKEN"))
			if err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}

			logger.Info("Updating Chart to use version %d.%d.%d", context.Major, context.Minor, context.Patch)
			if err := helmReleaser.Render(*context, dir); err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}

			// make a chart archive
			logger.Info("Validating chart in workspace")
			ch, err := chartutil.LoadDir(dir)
			if err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}
			logger.Info("Creating final archive of chart release")
			name, err := chartutil.Save(ch, wd)
			if err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}
			defer os.RemoveAll(name)

			// calculate the sha
			shasum, err := calculateSha256(name)
			if err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}

			// publish chart to scm host
			scm := scm.NewGitHubClient(os.Getenv("GITHUB_TOKEN"))
			downloadPath, err := scm.PublishRelease(*context, name)
			if err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}

			if helmReleaser.HelmRepo != nil {
				if helmReleaser.HelmRepo.ChartServer != nil {
					// publish chart to chartserver
					spec, err := chartserver.CreateChartVersionSpec(helmReleaser, context, "default", shasum, downloadPath)
					if err != nil {
						logger.Error(err)
						logger.Info("")
						os.Exit(1)
						return nil
					}

					if helmReleaser.HelmRepo.ChartServer.Filename != "" {
						if err = chartserver.WriteSpecFile(spec, helmReleaser.HelmRepo.ChartServer.Filename); err != nil {
							logger.Error(err)
							logger.Info("")
							os.Exit(1)
							return nil
						}
					}

				}
			}

			logger.Info("")

			return nil
		},
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cobra.OnInitialize(initConfig)

	cmd.AddCommand(InitCmd(out))

	cmd.Flags().Bool("snapshot", false, "when set, a snapshot release is created instead of a full release")
	cmd.Flags().String("chart-dir", wd, "path to chart")

	viper.BindPFlags(cmd.Flags())

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func InitAndExecute() {
	if err := RootCmd(os.Stdout).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetEnvPrefix("HELMRELEASER")
	viper.AutomaticEnv()
}

func calculateSha256(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", errors.Wrap(err, "failed to copy file to sha hasher")
	}

	return fmt.Sprintf("%#x", hasher.Sum(nil)), nil
}
