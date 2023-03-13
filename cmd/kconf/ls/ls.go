package ls

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sn3d/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "ls",
	Usage: "list of all contexts",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
		&cli.BoolFlag{
			Name:  "l",
			Usage: "use a long listing format",
			Value: false,
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		kubeConfigPath := cCtx.String("kubeconfig")
		longList := cCtx.Bool("l")

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

		w := tabwriter.NewWriter(os.Stdout, 5, 2, 2, ' ', tabwriter.DiscardEmptyColumns)
		if longList {
			fmt.Fprintf(w, "CONTEXT\tCLUSTER\tUSER\tNAMESPACE\n")
		}
		for _, ctx := range kc.Contexts {
			if longList {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ctx.Name, ctx.Context.Cluster, ctx.Context.User, ctx.Context.Namespace)
			} else {
				fmt.Fprintf(w, "%s\n", ctx.Name)
			}
		}
		w.Flush()

		return nil
	},
}
