package cc

import (
	"fmt"
	"os"
	"strings"

	"github.com/sn3d/kconf"
	"github.com/sn3d/kconf/internal/tui"
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

		var selected string
		if cCtx.Args().First() != "" {
			selected = cCtx.Args().First()
		} else {
			selected = showList(kubeConfigFile, kc)
		}

		if selected == "" {
			return fmt.Errorf("nothing selected\n")
		}

		kc.CurrentContext = selected
		err = kc.Save(kubeConfigFile)
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
