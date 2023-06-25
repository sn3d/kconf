package rm

import (
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "rm",
	Usage:     "remove context and context's cluster and user (if it's possible)",
	ArgsUsage: "[contextName]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		contextName := cCtx.Args().First()

		kc, path, err := kconf.Open(cCtx.String("kubeconfig"))
		if err != nil {
			return err
		}

		kc.Remove(contextName)

		err = kc.Save(path)
		if err != nil {
			return err
		}

		return nil
	},
}
