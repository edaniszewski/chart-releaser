package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCollector(t *testing.T) {
	c := NewCollector()
	assert.NotNil(t, c)
	assert.Empty(t, c.errors)
}

func TestCollector_Count(t *testing.T) {
	c := Collector{}
	assert.Equal(t, 0, c.Count())

	c.errors = append(c.errors, errors.New("test"))
	assert.Equal(t, 1, c.Count())
}

func TestCollector_HasErrors(t *testing.T) {
	c := Collector{}
	assert.False(t, c.HasErrors())

	c.errors = append(c.errors, errors.New("test"))
	assert.True(t, c.HasErrors())
}

func TestCollector_Add(t *testing.T) {
	c := Collector{}
	assert.Len(t, c.errors, 0)

	c.Add(errors.New("test 1"))
	assert.Len(t, c.errors, 1)

	c.Add(errors.New("test 2"))
	assert.Len(t, c.errors, 2)

	c.Add(errors.New("test 3"))
	assert.Len(t, c.errors, 3)
}

func TestCollector_Add2(t *testing.T) {
	c1 := Collector{}
	c2 := Collector{}

	c1.Add(errors.New("test  1"))
	assert.Len(t, c1.errors, 1)

	c1.Add(errors.New("test 2"))
	assert.Len(t, c1.errors, 2)

	c2.Add(errors.New("test 3"))
	assert.Len(t, c2.errors, 1)

	c2.Add(&c1)
	assert.Len(t, c2.errors, 3)
}

func TestCollector_Error_NoErrors(t *testing.T) {
	c := Collector{}
	assert.Equal(t, "", c.Error())
}

func TestCollector_Error_MultipleErrors(t *testing.T) {
	c := Collector{
		errors: []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		},
	}

	expected := `
Errors:
 • error 1
 • error 2
 • error 3

`

	assert.Equal(t, expected, c.Error())
}
