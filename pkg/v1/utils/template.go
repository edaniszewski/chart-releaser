package utils

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/apex/log"
	context "github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
)

// RenderTemplate is a convenience method to render a template, handing the case
// where dry-run is configured.
func RenderTemplate(ctx *context.Context, name, tmpl string) (string, error) {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		if err := ctx.CheckDryRun(err); err != nil {
			return "", err
		}
		t = template.Must(template.New("").Parse("dry-run"))
		log.Warnf("dry-run: failed to parse template '%s', using stand-in", name)
	}

	buf := bytes.Buffer{}
	err = t.Execute(&buf, ctx)
	if err != nil {
		if err := ctx.CheckDryRun(err); err != nil {
			return "", err
		}
		_, _ = fmt.Fprint(&buf, "dry-run")
		log.Warnf("dry-run: failed to execute template '%s', using stand-in", name)
	}
	return buf.String(), nil
}
