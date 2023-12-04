package mod

import (
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "mod",
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
		&cli.StringFlag{
			Name:  "client-certificate",
			Usage: "set the client-certificate from given file",
		},
		&cli.StringFlag{
			Name:  "client-key",
			Usage: "set the client-key from given file",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		kc, err := kconf.Open(cCtx.String("kubeconfig"))
		if err != nil {
			return err
		}

		opts := &kconf.UsermodOptions{
			Token:                 cCtx.String("token"),
			ClientCertificateFile: cCtx.String("client-certificate"),
			ClientKeyFile:         cCtx.String("client-key"),
		}

		err = kc.Usermod(cCtx.Args().First(), opts)
		if err != nil {
			return err
		}

		err = kc.Save()
		if err != nil {
			return err
		}

		return nil
	},
}
