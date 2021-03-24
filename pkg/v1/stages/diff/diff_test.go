package diff

import (
	"bytes"
	"testing"

	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "diff", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "displaying changes to chart files", Stage{}.String())
}

func TestStage_Run_NoDiff(t *testing.T) {
	context := ctx.Context{
		ShowDiff: false,
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
}

func TestStage_Run(t *testing.T) {
	buf := bytes.Buffer{}
	context := ctx.Context{
		ShowDiff: true,
		Out:      &buf,
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "foo/bar/Chart.yaml",
				PreviousContents: []byte("version: 1"),
				NewContents:      []byte("version: 2"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "testfile.txt",
				PreviousContents: []byte("abc"),
				NewContents:      []byte("123"),
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.Equal(t, "\x1b[0;37m+----start diff\n\x1b[0m\x1b[0;34m===\x1b[0m\nShowing changes to \x1b[0;33mChart.yaml\x1b[0m\n\n\x1b[0;31m- version: 1\x1b[0m\n\x1b[0;32m+ version: 2\x1b[0m\n\x1b[0;34m===\x1b[0m\nShowing changes to \x1b[0;33mtestfile.txt\x1b[0m\n\n\x1b[0;31m- abc\x1b[0m\n\x1b[0;32m+ 123\x1b[0m\n\x1b[0;37m+----end diff\n\x1b[0m", buf.String())
}
