package main

import (
	"fmt"
	"github.com/sn3d/kconf/cmd/kconf/cc"
	"os"

	"github.com/sn3d/kconf/cmd/kconf/chclus"
	"github.com/sn3d/kconf/cmd/kconf/chns"
	"github.com/sn3d/kconf/cmd/kconf/chusr"
	"github.com/sn3d/kconf/cmd/kconf/clustermod"
	"github.com/sn3d/kconf/cmd/kconf/export"
	"github.com/sn3d/kconf/cmd/kconf/imprt"
	"github.com/sn3d/kconf/cmd/kconf/ls"
	"github.com/sn3d/kconf/cmd/kconf/mv"
	"github.com/sn3d/kconf/cmd/kconf/rm"
	"github.com/sn3d/kconf/cmd/kconf/split"
	"github.com/sn3d/kconf/cmd/kconf/usermod"
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
			split.Cmd,
			cc.Cmd,
			chns.Cmd,
			chusr.Cmd,
			chclus.Cmd,
			usermod.Cmd,
			clustermod.Cmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		// print error
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
