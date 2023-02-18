package imprt

import (
	"io"
	"os"

	"github.com/sn3d/kconf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "import",
	Usage: "Import given kubeconfig on stdin to your configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig where to import",
		},
		&cli.BoolFlag{
			Name:  "base64",
			Usage: "Use if your input is base64 decoded kubeconfig",
		},
	},

	// main entry point for 'import'
	Action: func(cCtx *cli.Context) error {
		var kc *kconf.KubeConfig
		var err error

		kubeConfigFile := cCtx.String("kubeconfig")
		if kubeConfigFile != "" {
			kc, err = kconf.OpenFile(kubeConfigFile)
		} else {
			kc, err = kconf.OpenDefault()
		}

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
			sourceCfg, err = kconf.Open(data)
		}

		if err != nil {
			return err
		}

		kc.Import(sourceCfg)
		if kubeConfigFile != "" {
			err = kc.Save(kubeConfigFile)
		} else {
			err = kc.SaveDefault()
		}

		if err != nil {
			return err
		}

		return nil
	},
}
