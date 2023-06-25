package clustermod

import (
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
		kc, path, err := kconf.Open(cCtx.String("kubeconfig"))
		if err != nil {
			return err
		}

		opts := &kconf.ClustermodOptions{
			ServerURL: cCtx.String("url"),
		}

		err = kc.Clustermod(cCtx.Args().First(), opts)
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
