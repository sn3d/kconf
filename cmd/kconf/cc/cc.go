package cc

import (
	"fmt"
	"os"
	"strings"

	"github.com/sn3d/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "cc",
	Usage:     "change current context",
	ArgsUsage: "[CONTEXT]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		var kc *kconf.KubeConfig
		var err error

		kubeConfigFile := cCtx.String("kubeconfig")
		if kubeConfigFile == "" {
			configs := strings.Split(os.Getenv("KUBECONFIG"), ":")
			kubeConfigFile = configs[0]
		}

		kc, err = kconf.Open(kubeConfigFile)
		if err != nil {
			fmt.Printf("Cannot open your kubeconfig. Check if you have KUBECONFIG env. variable defined, or use --kubeconfig.\n")
		}

		kc.CurrentContext = cCtx.Args().First()

		err = kc.Save(kubeConfigFile)
		if err != nil {
			return err
		}

		return nil
	},
}
