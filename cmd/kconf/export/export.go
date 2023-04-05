package export

import (
	"os"

	"github.com/sn3d/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "export",
	Usage:     "Export context with user and cluster from your configuration",
	ArgsUsage: "[contextName]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to the kubeconfig file from where context is exported",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		contextName := cCtx.Args().First()
		kubeConfigPath := cCtx.String("kubeconfig")

		var kc *kconf.KubeConfig
		var err error

		kc, err = kconf.Open(kubeConfigPath)
		if err != nil {
			return err
		}

		exportedCfg, err := kc.Export(contextName)
		if err != nil {
			return err
		}

		exportedCfg.WriteTo(os.Stdout)
		return nil
	},
}
