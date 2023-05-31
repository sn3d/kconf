package chusr

import (
	"fmt"
	"os"
	"strings"

	"github.com/sn3d/kconf"
	"github.com/sn3d/kconf/internal/tui"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "chusr",
	Usage:     "change user for context",
	ArgsUsage: "[USER]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
		&cli.StringFlag{
			Name:    "context",
			Aliases: []string{"c"},
			Usage:   "context for which you want to change user",
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
			selected = showUserList(cCtx.String("context"), kc)
		}

		err = kc.Chusr(cCtx.String("context"), selected)
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

func showUserList(context string, conf *kconf.KubeConfig) string {
	if context == "" {
		context = conf.CurrentContext
	}

	opts := make([]string, len(conf.AuthInfos))
	for i := range conf.AuthInfos {
		opts[i] = conf.AuthInfos[i].Name
	}

	title := fmt.Sprintf("change user for '%s' context ", context)
	selected, _ := tui.ShowSimpleList(title, "", opts)
	return selected
}
