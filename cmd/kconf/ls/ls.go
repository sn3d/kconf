package ls

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sn3d/kconf/pkg/kconf"
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

		kc, err = kconf.Open(kubeConfigPath)

		if err != nil {
			fmt.Printf("Cannot open your kubeconfig. Check if you have KUBECONFIG env. variable defined, or use --kubeconfig.\n")
		}

		w := tabwriter.NewWriter(os.Stdout, 5, 2, 2, ' ', tabwriter.DiscardEmptyColumns)
		if longList {
			fmt.Fprintf(w, "CONTEXT\tCLUSTER\tUSER\tNAMESPACE\n")
		}
		for _, ctx := range kc.Contexts {
			if longList {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ctx.Name, ctx.Context.Cluster, ctx.Context.AuthInfo, ctx.Context.Namespace)
			} else {
				fmt.Fprintf(w, "%s\n", ctx.Name)
			}
		}
		w.Flush()

		return nil
	},
}
