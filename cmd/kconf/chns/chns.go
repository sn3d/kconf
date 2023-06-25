package chns

import (
	"fmt"

	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "chns",
	Usage:     "change default namespace for context",
	ArgsUsage: "[NAMESPACE]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
		&cli.StringFlag{
			Name:    "context",
			Aliases: []string{"c"},
			Usage:   "context for which you want to change namespace",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		kc, path, err := kconf.Open(cCtx.String("kubeconfig"))
		if err != nil {
			return err
			fmt.Printf("Cannot open your kubeconfig. Check if you have KUBECONFIG env. variable defined, or use --kubeconfig.\n")
		}

		err = kc.Chns(cCtx.String("context"), cCtx.Args().First())
		if err != nil {
			return err
		}

		err = kc.Save(path)
		if err != nil {
			return err
		}

		return nil
	},
}
