package mv

import (
	"github.com/sn3d/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "mv",
	Usage:     "rename SOURCE context to DEST context. Also rename context's cluster and user if it's possible",
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

		kubeConfigPath := cCtx.String("kubeconfig")

		var kc *kconf.KubeConfig
		var err error

		if kubeConfigPath == "" {
			kc, err = kconf.OpenDefault()
		} else {
			kc, err = kconf.OpenFile(kubeConfigPath)
		}

		if err != nil {
			return err
		}

		kc.RenameContext(src, dest)

		if kubeConfigPath == "" {
			err = kc.SaveDefault()
		} else {
			err = kc.Save(kubeConfigPath)
		}

		if err != nil {
			return err
		}

		return nil
	},
}
