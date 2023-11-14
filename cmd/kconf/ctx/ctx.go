package ctx

import (
	"fmt"
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/sn3d/kconf/pkg/tui"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "ctx",
	Usage:     "change current context",
	ArgsUsage: "[CONTEXT]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
	},

	Subcommands: []*cli.Command{
		mvCmd,
		rmCmd,
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
			selected = showList(path, kc)
		}

		if selected == "" {
			return fmt.Errorf("nothing selected")
		}

		kc.CurrentContext = selected
		err = kc.Save(path)
		if err != nil {
			return err
		}

		return nil
	},
}

func showList(file string, conf *kconf.KubeConfig) string {
	selected, _ := tui.ShowContextList(file, conf.CurrentContext, conf)
	return selected
}
