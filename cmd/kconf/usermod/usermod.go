package usermod

import (
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "usermod",
	Usage:     "modify a user",
	ArgsUsage: "[USER]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
		&cli.StringFlag{
			Name:  "token",
			Usage: "set the user token",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		kc, path, err := kconf.Open(cCtx.String("kubeconfig"))
		if err != nil {
			return err
		}

		opts := &kconf.UsermodOptions{
			Token: cCtx.String("token"),
		}

		err = kc.Usermod(cCtx.Args().First(), opts)
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
