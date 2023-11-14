package clst

import (
	"fmt"
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/sn3d/kconf/pkg/tui"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "clst",
	Usage:     "change cluster for context",
	ArgsUsage: "[CLUSTER]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
		&cli.StringFlag{
			Name:    "context",
			Aliases: []string{"c"},
			Usage:   "context for which you want to change cluster",
		},
	},

	Subcommands: []*cli.Command{
		modCmd,
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		kc, path, err := kconf.Open(cCtx.String("kubeconfig"))
		if err != nil {
			return err
		}

		var selected string
		if cCtx.Args().First() != "" {
			selected = cCtx.Args().First()
		} else {
			selected = showClusterList(cCtx.String("context"), kc)
		}

		err = kc.Chclus(cCtx.String("context"), selected)
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

func showClusterList(context string, conf *kconf.KubeConfig) string {
	if context == "" {
		context = conf.CurrentContext
	}

	opts := make([]string, len(conf.Clusters))
	for i := range conf.Clusters {
		opts[i] = conf.Clusters[i].Name
	}

	title := fmt.Sprintf("change cluster for '%s' context ", context)
	selected, _ := tui.ShowSimpleList(title, "", opts)
	return selected
}
