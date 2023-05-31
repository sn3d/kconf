package clustermod

import (
	"fmt"
	"os"
	"strings"

	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "clustermod",
	Usage:     "modify a cluster",
	ArgsUsage: "[CLUSTER]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
		&cli.StringFlag{
			Name:  "url",
			Usage: "set the cluster URL",
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

		opts := &kconf.ClustermodOptions{
			ServerURL: cCtx.String("url"),
		}

		err = kc.Clustermod(cCtx.Args().First(), opts)
		if err != nil {
			return err
		}

		err = kc.Save(kubeConfigFile)
		if err != nil {
			return err
		}

		return nil
	},
}
