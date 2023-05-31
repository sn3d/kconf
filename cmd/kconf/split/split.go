package split

import (
	"fmt"
	"os"
	"strings"

	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "split",
	Usage:     "split a kubeconfig into separated context pieces",
	ArgsUsage: "[PREFIX]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig you wish to split",
		},
		&cli.StringFlag{
			Name:        "additional-suffix",
			Usage:       "append an additional SUFFIX to file names",
			DefaultText: "",
		},
		&cli.BoolFlag{
			Name:    "numeric-suffixes",
			Aliases: []string{"d"},
			Usage:   "Use 2 digit numeric as suffix, not a context name",
		},
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
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

		splitted := kc.Split()

		for i := range splitted {
			fileName := createFileName(i, splitted[i], cCtx)
			splitted[i].Save(fileName)
			fmt.Fprintf(os.Stdout, "%s\n", fileName)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func createFileName(idx int, kcfg *kconf.KubeConfig, cCtx *cli.Context) string {
	prefix := cCtx.Args().First()
	additionalSuffix := cCtx.String("additional-suffix")

	suffix := kcfg.CurrentContext
	if cCtx.Bool("numeric-suffixes") {
		suffix = fmt.Sprintf("%02d", idx)
	}

	return fmt.Sprintf("%s%s%s", prefix, suffix, additionalSuffix)
}
