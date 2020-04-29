package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/config"
	v1 "github.com/edaniszewski/chart-releaser/pkg/v1"
	"github.com/spf13/cobra"
)

type checkCmd struct {
	c *cobra.Command

	skipEnv bool
}

func newCheckCommand() *checkCmd {
	root := &checkCmd{}
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check if a configuration file is valid",
		Long: heredoc.Doc(`
			This command checks that the specified file is a valid chart-releaser
			configuration. If the provided path is a directory, this will look for
			.chartreleaser.yml within that directory.
	
			If no path is specified, this will look for .chartreleaser.yml in the
			current working directory.
	
			This will error if no file is found, the file can't be loaded, required
			fields are missing, or if any expected environment variables are not set.
			The --skip-env flag can be set to skip the checks for environment variables.
		`),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.WithFields(log.Fields{
				"cmd": "check",
			}).Debug("running command")

			var p string
			if len(args) == 1 {
				p = args[0]
			}

			v, err := config.Load(p)
			if err != nil {
				return err
			}

			switch v.GetVersion() {
			case v1.ConfigVersion():
				err = v1.NewChecker(v.GetData()).Run(v1.CheckerOptions{
					SkipEnv: root.skipEnv,
				})
			default:
				err = fmt.Errorf("unsupported config version: %s", v.GetVersion())
			}

			if err == nil {
				fmt.Printf("successfully loaded %v config file (%v)\n", v.GetVersion(), v.GetPath())
			}
			return err
		},
	}

	cmd.Flags().BoolVar(&root.skipEnv, "skip-env", false, "skip checks for expected environment variables")

	root.c = cmd
	return root
}
