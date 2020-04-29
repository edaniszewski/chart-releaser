package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintDiff_EmptyStrings(t *testing.T) {
	buf := bytes.Buffer{}

	PrintDiff(&buf, "testfile.txt", "", "")
	assert.Equal(
		t,
		"\x1b[34m===\x1b[0m\nShowing changes to \x1b[33mtestfile.txt\x1b[0m\n\n\x1b[33mno changes\x1b[0m\n",
		buf.String(),
	)
}

func TestPrintDiff_NoChangeSingleLine(t *testing.T) {
	buf := bytes.Buffer{}

	PrintDiff(&buf, "testfile.txt", "abc", "abc")
	assert.Equal(
		t,
		"\x1b[34m===\x1b[0m\nShowing changes to \x1b[33mtestfile.txt\x1b[0m\n\n\x1b[33mno changes\x1b[0m\n",
		buf.String(),
	)
}

func TestPrintDiff_NoChangeMultiLine(t *testing.T) {
	buf := bytes.Buffer{}

	PrintDiff(&buf, "testfile.txt", "a\nb\nc\n", "a\nb\nc\n")
	assert.Equal(
		t,
		"\x1b[34m===\x1b[0m\nShowing changes to \x1b[33mtestfile.txt\x1b[0m\n\n\x1b[33mno changes\x1b[0m\n",
		buf.String(),
	)
}

func TestPrintDiff_SingleLine(t *testing.T) {
	buf := bytes.Buffer{}

	PrintDiff(&buf, "testfile.txt", "abc", "acb")
	assert.Equal(
		t,
		"\x1b[34m===\x1b[0m\nShowing changes to \x1b[33mtestfile.txt\x1b[0m\n\n\x1b[31m- abc\x1b[0m\n\x1b[32m+ acb\x1b[0m\n",
		buf.String(),
	)
}

func TestPrintDiff_MultiLine(t *testing.T) {
	buf := bytes.Buffer{}

	PrintDiff(&buf, "testfile.txt", "a\nb\nc\n", "a\nc\nb\n")
	assert.Equal(
		t,
		"\x1b[34m===\x1b[0m\nShowing changes to \x1b[33mtestfile.txt\x1b[0m\n\n  a\n\x1b[31m- b\x1b[0m\n  c\n\x1b[32m+ b\x1b[0m\n  \n",
		buf.String(),
	)
}
