package cli

import (
	"io"
	"os"
	"path"

	"github.com/helmreleaser/helmreleaser/pkg/helmreleaser"
	"github.com/helmreleaser/helmreleaser/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "init",
		Short:         "Generates a new .helmreleaser.yaml file",
		Long:          ``,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logger.NewLogger()

			logger.Info("")
			defer func() {
				logger.Info("")
			}()

			wd, err := os.Getwd()
			if err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}
			helmreleaser := helmreleaser.CreateDefault()
			logger.Info("Creating new .helmreleaser.yaml file")
			p := path.Join(wd, ".helmreleaser.yaml")
			if err := helmreleaser.WriteToFile(p, false); err != nil {
				logger.Error(err)
				logger.Info("")
				os.Exit(1)
				return nil
			}

			logger.Info("Default config created at %s, please edit", p)
			return nil
		},
	}

	return cmd
}
