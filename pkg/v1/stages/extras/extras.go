package extras

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"

	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/client"
	context "github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
)

// Stage for the "extras" step of the update pipeline.
type Stage struct{}

// Name of the stage.
func (Stage) Name() string {
	return "extras"
}

// String describes what the stage does.
func (Stage) String() string {
	return "updating additional chart files"
}

// Run the operations defined for the stage.
func (Stage) Run(ctx *context.Context) error {

	for _, extra := range ctx.Config.Extras {

		extraFile := context.File{
			Path: extra.Path,
		}

		opts := &client.Options{
			RepoName:  ctx.Repository.Name,
			RepoOwner: ctx.Repository.Owner,
		}
		log.WithFields(log.Fields{
			"path":      extra.Path,
			"repoName":  ctx.Repository.Name,
			"repoOwner": ctx.Repository.Owner,
		}).Debug("getting file contents")

		contents, err := ctx.Client.GetFile(ctx.Context, opts, extra.Path)
		if err != nil {
			if err := ctx.CheckDryRun(err); err != nil {
				return err
			}
			log.WithField("path", extra.Path).Warn("failed to get contents for file -- skipping")
			continue
		}
		extraFile.PreviousContents = []byte(contents)

		for _, update := range extra.Updates {
			re, err := regexp.Compile(update.Search)
			if err != nil {
				if err := ctx.CheckDryRun(err); err != nil {
					return err
				}
				log.WithField("re", re).Warn("failed to compile search regex")
				continue
			}

			// Get the limit of splits. If the zero value is set, default to -1
			// (split on all occurrences), as  a limit of 0 means do not split.
			limit := update.Limit
			if limit == 0 {
				limit = -1
			}
			parts := re.Split(contents, limit)

			// Parse the replace string as a template. If it is not a template, it will
			// remain unchanged.
			t, err := template.New("").Parse(update.Replace)
			if err != nil {
				if err := ctx.CheckDryRun(err); err != nil {
					return err
				}
				log.WithField("template", update.Replace).Warn("failed to parse string as template")
				continue
			}

			buf := bytes.Buffer{}
			err = t.Execute(&buf, ctx)
			if err != nil {
				if err := ctx.CheckDryRun(err); err != nil {
					return err
				}
				log.WithField("template", update.Replace).Warn("failed to execute template")
				continue
			}
			contents = strings.Join(parts, buf.String())
		}

		extraFile.NewContents = []byte(contents)

		if !extraFile.HasChanges() {
			log.WithFields(log.Fields{
				"contents":  string(extraFile.NewContents),
				"path":      extra.Path,
				"repoName":  ctx.Repository.Name,
				"repoOwner": ctx.Repository.Owner,
			}).Warn("no change detected to extras file")
		}

		ctx.Files = append(ctx.Files, extraFile)
	}
	return nil
}
