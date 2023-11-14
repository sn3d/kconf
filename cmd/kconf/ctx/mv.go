package ctx

import (
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/urfave/cli/v2"
)

var mvCmd = &cli.Command{
	Name:      "mv",
	Usage:     "move source context to destination context",
	ArgsUsage: "[SOURCE] [DEST]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		src := cCtx.Args().Get(0)
		dest := cCtx.Args().Get(1)

		kc, path, err := kconf.Open(cCtx.String("kubeconfig"))
		if err != nil {
			return err
		}

		kc.Rename(src, dest)

		err = kc.Save(path)
		if err != nil {
			return err
		}

		return nil
	},
}
