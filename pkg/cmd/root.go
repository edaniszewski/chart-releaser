package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

// Execute is the entry point for the command line tool. It runs the root command.
func Execute(exiter func(int), args []string) {
	root := newRootCommand(exiter)
	root.SetArgs(args)

	if err := root.Execute(); err != nil {
		root.exiter(1)
	}
}

type rootCmd struct {
	c *cobra.Command

	exiter func(int)
	debug  bool
}

func newRootCommand(exiter func(int)) *rootCmd {
	root := &rootCmd{
		exiter: exiter,
	}

	cmd := &cobra.Command{
		Use:   "chart-releaser",
		Short: "Update Helm Chart versions for new application releases",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if root.debug {
				log.SetLevel(log.DebugLevel)
				log.Debug("debug logging enabled")
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&root.debug, "debug", false, "run chart-releaser with debug logging")

	cmd.AddCommand(
		newCheckCommand().c,
		newFmtCommand().c,
		newInitCommand().c,
		newUpdateCommand().c,
		newVersionCommand().c,
	)

	root.c = cmd
	return root
}

// SetArgs is an alias to the underlying cobra.Command's SetArgs.
func (c *rootCmd) SetArgs(a []string) {
	c.c.SetArgs(a)
}

// Execute is an alias to the underlying cobra.Command's Execute.
func (c *rootCmd) Execute() error {
	return c.c.Execute()
}
