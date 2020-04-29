package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/edaniszewski/chart-releaser/pkg/config"
	v1 "github.com/edaniszewski/chart-releaser/pkg/v1"
	"github.com/spf13/cobra"
)

type fmtCmd struct {
	c *cobra.Command

	noHeader bool
}

func newFmtCommand() *fmtCmd {
	root := &fmtCmd{}
	cmd := &cobra.Command{
		Use:   "fmt",
		Short: "Format a configuration file",
		Long: heredoc.Doc(`
			This command allows you to format the chart-releaser configuration file so
			that it follows a standard layout. This is entirely optional and provided
			as a convenience.
	
			Note that formatting removes all inline comments.
	
			The header comment is included by default but can be excluded with the
			--no-header flag.
		`),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var p string
			if len(args) == 1 {
				p = args[0]
			}

			path, err := config.GetConfigPath(p)
			if err != nil {
				return err
			}

			v, err := config.Load(p)
			if err != nil {
				return err
			}

			switch v.GetVersion() {
			case v1.ConfigVersion():
				err = v1.NewFormatter(v.GetData()).Run(v1.FormatterOptions{
					NoHeader: root.noHeader,
					Path:     path,
				})
			default:
				err = fmt.Errorf("unsupported config version: %s", v.GetVersion())
			}

			return err
		},
	}

	cmd.Flags().BoolVar(&root.noHeader, "no-header", false, "exclude the header comment when formatting")

	root.c = cmd
	return root
}
