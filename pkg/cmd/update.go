package cmd

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/config"
	v1 "github.com/edaniszewski/chart-releaser/pkg/v1"
	"github.com/spf13/cobra"
)

type updateCmd struct {
	c *cobra.Command

	timeout    time.Duration
	dryRun     bool
	allowDirty bool
	diff       bool
}

func newUpdateCommand() *updateCmd {
	root := &updateCmd{}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the Helm Chart for a new project release",
		Long: heredoc.Doc(`
			The project may be configured via environment, config file, or command line
			arguments.
		`),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var p string
			if len(args) == 1 {
				p = args[0]
			}
			start := time.Now()
			log.Info("starting chart update...")

			v, err := config.Load(p)
			if err != nil {
				return err
			}

			switch v.GetVersion() {
			case v1.ConfigVersion():
				err = v1.NewUpdater(v.GetData()).Run(v1.UpdateOptions{
					AllowDirty: root.allowDirty,
					DryRun:     root.dryRun,
					ShowDiff:   root.diff,
					Timeout:    root.timeout,
				})
			default:
				err = fmt.Errorf("unsupported config version: %s", v.GetVersion())
			}

			if err == nil {
				log.Debugf("update completed after %0.3fs", time.Since(start).Seconds())
			}
			return err
		},
	}

	cmd.Flags().BoolVar(&root.dryRun, "dry-run", false, "run the command without side effects")
	cmd.Flags().BoolVar(&root.allowDirty, "allow-dirty", false, "do not fail if the git repo is in a dirty state")
	cmd.Flags().BoolVar(&root.diff, "diff", false, "show the diff for files modified by the command")
	cmd.Flags().DurationVar(&root.timeout, "timeout", 5*time.Minute, "timeout for the entire update process")

	root.c = cmd
	return root
}
