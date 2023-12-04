package main

import (
	"fmt"
	"os"

	"github.com/sn3d/kconf/cmd/kconf/cluster"
	"github.com/sn3d/kconf/cmd/kconf/context"
	"github.com/sn3d/kconf/cmd/kconf/ns"
	"github.com/sn3d/kconf/cmd/kconf/user"

	"github.com/sn3d/kconf/cmd/kconf/export"
	"github.com/sn3d/kconf/cmd/kconf/imprt"
	"github.com/sn3d/kconf/cmd/kconf/split"
	"github.com/urfave/cli/v2"
)

// version is set by goreleaser, via -ldflags="-X 'main.version=...'".
var version = "development"

func main() {
	app := &cli.App{
		Name:                 "kconf",
		Version:              version,
		Usage:                "managing your kubeconfig",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			imprt.Cmd,
			export.Cmd,
			context.Cmd,
			user.Cmd,
			cluster.Cmd,
			ns.Cmd,
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
