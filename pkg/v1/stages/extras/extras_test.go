package extras

import (
	"errors"
	"testing"

	"github.com/edaniszewski/chart-releaser/internal/testutils"
	v1 "github.com/edaniszewski/chart-releaser/pkg/v1/cfg"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "extras", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "updating additional chart files", Stage{}.String())
}

func TestStage_RunNoExtras(t *testing.T) {
	context := ctx.Context{
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.Len(t, context.Files, 0)
}

func TestStage_Run(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			FileData: "test file data search",
		},
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{
				{
					Path: "path1",
					Updates: []*v1.SearchReplace{
						{
							Search:  "search",
							Replace: "replace",
						},
					},
				},
				{
					Path: "path2",
					Updates: []*v1.SearchReplace{
						{
							Search:  "search",
							Replace: "replace",
							Limit:   2,
						},
					},
				},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.Len(t, context.Files, 2)

	f1 := context.Files[0]
	assert.Equal(t, "path1", f1.Path)
	assert.Equal(t, "test file data search", string(f1.PreviousContents))
	assert.Equal(t, "test file data replace", string(f1.NewContents))

	f2 := context.Files[1]
	assert.Equal(t, "path2", f2.Path)
	assert.Equal(t, "test file data search", string(f2.PreviousContents))
	assert.Equal(t, "test file data replace", string(f2.NewContents))
}

func TestStage_RunGetFileError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			GetFileError: []error{
				errors.New("test error"),
			},
		},
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{
				{
					Path: "path1",
					Updates: []*v1.SearchReplace{
						{
							Search:  "search",
							Replace: "replace",
						},
					},
				},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "test error")
	assert.Len(t, context.Files, 0)
}

func TestStage_RunGetFileError_DryRun(t *testing.T) {
	context := ctx.Context{
		DryRun: true,
		Client: &testutils.FakeClient{
			GetFileError: []error{
				errors.New("test error"),
			},
		},
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{
				{
					Path: "path1",
					Updates: []*v1.SearchReplace{
						{
							Search:  "search",
							Replace: "replace",
						},
					},
				},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.Len(t, context.Files, 0)
	assert.EqualError(t, context.Errors(), "\nErrors:\n • test error\n\n")
}

func TestStage_RunCompileSearchError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			FileData: "test file data",
		},
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{
				{
					Path: "path1",
					Updates: []*v1.SearchReplace{
						{
							Search:  "*",
							Replace: "replace",
						},
					},
				},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "error parsing regexp: missing argument to repetition operator: `*`")
	assert.Len(t, context.Files, 0)
}

func TestStage_RunCompileSearchError_DryRun(t *testing.T) {
	context := ctx.Context{
		DryRun: true,
		Client: &testutils.FakeClient{
			FileData: "test file data",
		},
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{
				{
					Path: "path1",
					Updates: []*v1.SearchReplace{
						{
							Search:  "*",
							Replace: "replace",
						},
					},
				},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.Len(t, context.Files, 1)
	assert.EqualError(t, context.Errors(), "\nErrors:\n • error parsing regexp: missing argument to repetition operator: `*`\n\n")
}

func TestStage_RunParseReplaceTemplateError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			FileData: "test file data",
		},
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{
				{
					Path: "path1",
					Updates: []*v1.SearchReplace{
						{
							Search:  "foo",
							Replace: "{{end}}",
						},
					},
				},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "template: :1: unexpected {{end}}")
	assert.Len(t, context.Files, 0)
}

func TestStage_RunParseReplaceTemplateError_DryRun(t *testing.T) {
	context := ctx.Context{
		DryRun: true,
		Client: &testutils.FakeClient{
			FileData: "test file data",
		},
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{
				{
					Path: "path1",
					Updates: []*v1.SearchReplace{
						{
							Search:  "foo",
							Replace: "{{end}}",
						},
					},
				},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.Len(t, context.Files, 1)
	assert.EqualError(t, context.Errors(), "\nErrors:\n • template: :1: unexpected {{end}}\n\n")
}

func TestStage_RunExecuteTemplateError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			FileData: "test file data",
		},
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{
				{
					Path: "path1",
					Updates: []*v1.SearchReplace{
						{
							Search:  "foo",
							Replace: "{{ .Does.Not.Exist }}",
						},
					},
				},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "template: :1:8: executing \"\" at <.Does.Not.Exist>: can't evaluate field Does in type *ctx.Context")
	assert.Len(t, context.Files, 0)
}

func TestStage_RunExecuteTemplateError_DryRun(t *testing.T) {
	context := ctx.Context{
		DryRun: true,
		Client: &testutils.FakeClient{
			FileData: "test file data",
		},
		Config: &v1.Config{
			Extras: []*v1.ExtrasConfig{
				{
					Path: "path1",
					Updates: []*v1.SearchReplace{
						{
							Search:  "foo",
							Replace: "{{ .Does.Not.Exist }}",
						},
					},
				},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.Len(t, context.Files, 1)
	assert.EqualError(t, context.Errors(), "\nErrors:\n • template: :1:8: executing \"\" at <.Does.Not.Exist>: can't evaluate field Does in type *ctx.Context\n\n")
}
