package main

import (
	"fmt"
	"os"

	"github.com/sn3d/kconf/cmd/kconf/export"
	"github.com/sn3d/kconf/cmd/kconf/imprt"
	"github.com/sn3d/kconf/cmd/kconf/ls"
	"github.com/sn3d/kconf/cmd/kconf/mv"
	"github.com/sn3d/kconf/cmd/kconf/rm"
	"github.com/urfave/cli/v2"
)

// version is set by goreleaser, via -ldflags="-X 'main.version=...'".
var version = "development"

func main() {
	app := &cli.App{
		Name:    "kconf",
		Version: version,
		Usage:   "managing your kubeconfig",
		Commands: []*cli.Command{
			imprt.Cmd,
			export.Cmd,
			rm.Cmd,
			mv.Cmd,
			ls.Cmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		// print error
		fmt.Printf("Error: %e", err)
	}
}
