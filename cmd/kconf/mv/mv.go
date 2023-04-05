package mv

import (
	"os"
	"strings"

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

		kubeConfigFile := cCtx.String("kubeconfig")
		if kubeConfigFile == "" {
			configs := strings.Split(os.Getenv("KUBECONFIG"), ":")
			kubeConfigFile = configs[0]
		}

		var kc *kconf.KubeConfig
		var err error

		kc, err = kconf.Open(kubeConfigFile)
		if err != nil {
			return err
		}

		kc.Rename(src, dest)

		err = kc.Save(kubeConfigFile)
		if err != nil {
			return err
		}

		return nil
	},
}
