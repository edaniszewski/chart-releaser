package cmd

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	"github.com/edaniszewski/chart-releaser/pkg"
	"github.com/edaniszewski/chart-releaser/pkg/templates"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-latest"
)

type versionCmd struct {
	c *cobra.Command

	checkLatest bool
}

func newVersionCommand() *versionCmd {
	root := &versionCmd{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version info of chart-releaser",
		Long: heredoc.Doc(`
			This command prints out version information, including info about when
			the binary was built, for which platform, etc.

			If the --latest flag is set, this will also print out the latest release
			version from GitHub. You can compare this with your installed version to
			see if you can update your installed binary.
		`),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var versionBuffer = new(bytes.Buffer)
			var latestBuffer = new(bytes.Buffer)

			t := template.Must(template.New("version").Parse(templates.CommandVersionTemplate))
			if err := t.Execute(versionBuffer, pkg.NewVersionInfo()); err != nil {
				return err
			}

			if root.checkLatest {
				tag := &latest.GithubTag{
					Owner:      "edaniszewski",
					Repository: "chart-releaser",
				}

				res, err := latest.Check(tag, pkg.Version)
				if err != nil {
					return err
				}

				var status string
				if res.Outdated {
					status = "A new version of chart-releaser is available"
				} else if res.Latest {
					status = "Installed version is latest"
				}

				vs := pkg.LatestVersion{
					Status:    status,
					Latest:    res.Current,
					Installed: pkg.Version,
				}

				t := template.Must(template.New("latest-version").Parse(templates.FlagLatestVersionTemplate))
				if err := t.Execute(latestBuffer, &vs); err != nil {
					return err
				}
			}

			_, _ = fmt.Fprintln(os.Stdout, versionBuffer.String())
			if root.checkLatest {
				_, _ = fmt.Fprintln(os.Stdout, latestBuffer.String())
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&root.checkLatest, "latest", false, "print the latest releases version")

	root.c = cmd
	return root
}
