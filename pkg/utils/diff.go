package utils

import (
	"fmt"
	"io"
	"strings"

	"github.com/aryann/difflib"
	"github.com/mgutz/ansi"
)

// PrintDiff writes a git-like diff to the supplied io.Writer.
func PrintDiff(out io.Writer, filename, old, new string) {
	fmt.Fprintf(out, "%s\n", ansi.Color("===", "blue"))
	fmt.Fprintf(out, "Showing changes to %s\n\n", ansi.Color(filename, "yellow"))

	if old == new {
		fmt.Fprintf(out, "%s\n", ansi.Color("no changes", "yellow"))
		return
	}

	records := difflib.Diff(
		strings.Split(old, "\n"),
		strings.Split(new, "\n"),
	)
	for _, diff := range records {
		text := diff.Payload

		switch diff.Delta {
		case difflib.RightOnly:
			fmt.Fprintf(out, "%s\n", ansi.Color("+ "+text, "green"))
		case difflib.LeftOnly:
			fmt.Fprintf(out, "%s\n", ansi.Color("- "+text, "red"))
		case difflib.Common:
			fmt.Fprintf(out, "%s\n", "  "+text)
		}
	}
}
