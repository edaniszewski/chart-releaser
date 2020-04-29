package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/config"
	"github.com/edaniszewski/chart-releaser/pkg/templates"
	"github.com/edaniszewski/chart-releaser/pkg/utils"
	"github.com/spf13/cobra"
)

type initCmd struct {
	c *cobra.Command

	chart       string
	path        string
	githubOwner string
	githubRepo  string
	author      string
	email       string
}

func newInitCommand() *initCmd {
	root := &initCmd{}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate a new .chartreleaser.yaml config file",
		Long: heredoc.Doc(`
			If no directory is specified, the file will be created in the current
			working directory.
		`),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no path is provided, assume the current working directory.
			var rootPath = "."
			if len(args) == 1 {
				rootPath = args[0]
			}

			path, err := filepath.Abs(filepath.Join(rootPath, config.DefaultFile))
			if err != nil {
				return err
			}

			// If the config file already exists, there is nothing to init.
			if _, err := os.Stat(path); err == nil {
				return config.ErrConfigExists
			}

			chartName := root.chart
			if chartName == "" {
				p, err := filepath.Abs(rootPath)
				if err != nil {
					return err
				}
				_, chartName = filepath.Split(p)
			}

			opts, err := NewInitOptions(
				chartName,
				root.path,
				root.githubOwner,
				root.githubRepo,
				root.author,
				root.email,
			)

			tmpl, err := template.New("init").Parse(templates.ConfigFileTemplate)
			if err != nil {
				return err
			}
			var out bytes.Buffer
			if err := tmpl.Execute(&out, opts); err != nil {
				return err
			}

			return ioutil.WriteFile(path, out.Bytes(), 0644)
		},
	}

	cmd.Flags().StringVar(&root.chart, "chart", "", "the name  of the chart directory")
	cmd.Flags().StringVar(&root.path, "path", "", "the subdirectory within the repo containing the chart directory")
	cmd.Flags().StringVar(&root.githubOwner, "github-owner", "", "the user/organization for the GitHub repository containing the chart")
	cmd.Flags().StringVar(&root.githubRepo, "github-repo", "", "the name of the GitHub repository containing the chart")
	cmd.Flags().StringVar(&root.author, "author", "", "the name of the commit author")
	cmd.Flags().StringVar(&root.email, "email", "", "the email of the commit author")

	root.c = cmd
	return root
}

// InitOptions hold data for initializing a new chart-releaser config.
type InitOptions struct {
	Chart       string
	Path        string
	GithubOwner string
	GithubName  string
	AuthorName  string
	AuthorEmail string
	Header      string
}

// NewInitOptions creates a new InitOptions populated with the provided
// data, falling back to defaults where applicable.
func NewInitOptions(chart, path, ghOwner, ghName, author, email string) (*InitOptions, error) {
	if chart == "" {
		return nil, errors.New("no chart specified")
	}

	if ghOwner == "" {
		return nil, errors.New("no github owner specified")
	}

	// TODO: this could just be a string that we parse to figure out
	if ghName == "" {
		return nil, errors.New("no github name specified")
	}

	if author == "" {
		log.Debug("no author specified, checking git config")
		s, err := utils.GetGitUserName()
		if err != nil {
			return nil, err
		}
		author = s
	}

	if email == "" {
		log.Debug("no email specified, checking git config")
		s, err := utils.GetGitUserEmail()
		if err != nil {
			return nil, err
		}
		email = s
	}

	return &InitOptions{
		Chart:       chart,
		Path:        path,
		GithubOwner: ghOwner,
		GithubName:  ghName,
		AuthorEmail: email,
		AuthorName:  author,
		Header:      templates.ConfigHeaderComment,
	}, nil
}
