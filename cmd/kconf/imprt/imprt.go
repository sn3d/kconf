package imprt

import (
	"io"
	"os"

	"github.com/sn3d/kconf/pkg/kconf"
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
		&cli.StringFlag{
			Name:  "as",
			Usage: "import context, user and cluster AS. this option rename the imported context (only if it's one)",
		},
		&cli.BoolFlag{
			Name:  "base64",
			Usage: "if your input is base64 decoded kubeconfig",
		},
	},

	// main entry point for 'import'
	Action: func(cCtx *cli.Context) error {
		kc, path, err := kconf.Open(cCtx.String("kubeconfig"))
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

		opts := kconf.ImportOptions{As: cCtx.String("as")}
		kc.Import(sourceCfg, &opts)

		err = kc.Save(path)
		if err != nil {
			return err
		}

		return nil
	},
}
