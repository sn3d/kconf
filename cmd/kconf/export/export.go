package export

import (
	"fmt"

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
			Usage: "path to kubeconfig from where export context",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		contextName := cCtx.Args().First()
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

		exportedCfg, err := kc.Export(contextName)
		if err != nil {
			return err
		}

		fmt.Printf("%s", exportedCfg.ToString())
		return nil
	},
}
