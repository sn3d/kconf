package imprt

import (
	"io"
	"os"
	"strings"

	"github.com/sn3d/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "import",
	Usage: "import given kubeconfig on stdin to your configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to the dest. kubeconfig where context is imported",
		},
		&cli.BoolFlag{
			Name:  "base64",
			Usage: "if your input is base64 decoded kubeconfig",
		},
	},

	// main entry point for 'import'
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
			return err
		}

		// get configuration from stdin
		// it support also heredoc:
		//   kconf put << EOF
		//   Multi-line
		//   heredoc
		//   EOF
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		var sourceCfg *kconf.KubeConfig
		if cCtx.Bool("base64") {
			sourceCfg, err = kconf.OpenBase64(data)
		} else {
			sourceCfg, err = kconf.OpenData(data)
		}

		if err != nil {
			return err
		}

		kc.Import(sourceCfg)

		err = kc.Save(kubeConfigFile)
		if err != nil {
			return err
		}

		return nil
	},
}
