package main

import (
	"fmt"
	"github.com/sn3d/kconf/cmd/kconf/clst"
	"github.com/sn3d/kconf/cmd/kconf/ctx"
	"github.com/sn3d/kconf/cmd/kconf/ns"
	"github.com/sn3d/kconf/cmd/kconf/usr"
	"os"

	"github.com/sn3d/kconf/cmd/kconf/export"
	"github.com/sn3d/kconf/cmd/kconf/imprt"
	"github.com/sn3d/kconf/cmd/kconf/ls"
	"github.com/sn3d/kconf/cmd/kconf/split"
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
			ctx.Cmd,
			usr.Cmd,
			clst.Cmd,
			ns.Cmd,
			ls.Cmd,
			split.Cmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		// print error
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
