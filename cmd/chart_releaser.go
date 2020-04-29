package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/edaniszewski/chart-releaser/pkg/cmd"
)

func main() {
	log.SetHandler(cli.Default)

	cmd.Execute(
		os.Exit,
		os.Args[1:],
	)
}
