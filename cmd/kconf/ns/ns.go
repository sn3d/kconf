package ns

import (
	"fmt"
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/sn3d/kconf/pkg/tui"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "ns",
	Usage:     "change default namespace for context",
	ArgsUsage: "[NAMESPACE]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
		&cli.StringFlag{
			Name:    "context",
			Aliases: []string{"c"},
			Usage:   "context for which you want to change namespace",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		kc, path, err := kconf.Open(cCtx.String("kubeconfig"))
		if err != nil {
			return err
			fmt.Printf("Cannot open your kubeconfig. Check if you have KUBECONFIG env. variable defined, or use --kubeconfig.\n")
		}

		namespace := cCtx.Args().First()
		if namespace == "" {
			namespace = showNamespaceList(cCtx.String("context"), kc)
		}

		// nothing to change
		if namespace == "" {
			return nil
		}

		err = kc.Chns(cCtx.String("context"), namespace)
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

func showNamespaceList(configContext string, kc *kconf.KubeConfig) string {

	if configContext != "" {
		kc.CurrentContext = configContext
	}

	namespaces, err := kc.GetAllNamespaces()
	if err != nil {
		return ""
	}

	namespace, err := tui.ShowSimpleList(kc.CurrentContext+" namespaces", "", namespaces)
	if err != nil {
		return ""
	}

	return namespace
}
