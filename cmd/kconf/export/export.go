package export

import (
	"os"

	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "export",
	Usage:     "export context with user and cluster from your configuration",
	ArgsUsage: "[contextName]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to the kubeconfig file from where context is exported",
		},
		&cli.StringFlag{
			Name:  "as",
			Usage: "export context, user and cluster AS. this option rename the imported context (only if it's one)",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		contextName := cCtx.Args().First()

		kc, err := kconf.Open(cCtx.String("kubeconfig"))
		if err != nil {
			return err
		}

		opts := kconf.ExportOptions{
			As: cCtx.String("as"),
		}

		exportedCfg, err := kc.Export(contextName, &opts)
		if err != nil {
			return err
		}

		exportedCfg.WriteTo(os.Stdout)
		return nil
	},
}
