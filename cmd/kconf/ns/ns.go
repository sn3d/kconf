package ns

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/sn3d/kconf/pkg/tui/list"
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
		namespace := cCtx.Args().First()
		if namespace != "" {
			return directChange(cCtx, namespace)
		} else {
			return showTUI(cCtx)
		}
	},
}

func directChange(cCtx *cli.Context, namespace string) error {

	kc, err := kconf.Open(cCtx.String("kubeconfig"))
	if err != nil {
		return err
	}

	err = kc.ChangeNamespace(cCtx.String("context"), namespace)
	if err != nil {
		return err
	}

	return kc.Save()
}

func showTUI(cCtx *cli.Context) error {
	kc, err := kconf.Open(cCtx.String("kubeconfig"))
	if err != nil {
		return err
	}

	model, err := NewModel(kc, cCtx.String("context"))
	if err != nil {
		return err
	}

	m, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}

	out := m.(Model)
	switch msg := out.ExitMsg.(type) {
	case list.CloseMsg:
		if !msg.Aborted {
			fmt.Printf("changes saved\n")
		}
	case list.PickedMsg:
		fmt.Printf("namespace changed to %s\n", msg.Picked)
	}

	return err
}
