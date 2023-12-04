package context

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sn3d/kconf/cmd/kconf/context/ls"
	"github.com/sn3d/kconf/cmd/kconf/context/mv"
	"github.com/sn3d/kconf/cmd/kconf/context/rm"
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/sn3d/kconf/pkg/tui/list"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "ctx",
	Usage:     "context commands",
	ArgsUsage: "[CONTEXT]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
	},

	Subcommands: []*cli.Command{
		mv.Cmd,
		rm.Cmd,
		ls.Cmd,
	},

	Action: func(cCtx *cli.Context) error {
		k8sContext := cCtx.Args().First()
		if k8sContext != "" {
			return directChange(cCtx, k8sContext)
		} else {
			return showTUI(cCtx)
		}
	},
}

func directChange(cCtx *cli.Context, k8sContext string) error {

	kc, err := kconf.Open(cCtx.String("kubeconfig"))
	if err != nil {
		return err
	}

	kc.CurrentContext = k8sContext

	err = kc.Save()
	if err != nil {
		return err
	}

	return nil
}

func showTUI(cCtx *cli.Context) error {
	kc, err := kconf.Open(cCtx.String("kubeconfig"))
	if err != nil {
		return err
	}

	model, err := NewModel(kc)
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
		fmt.Printf("context changed to %s\n", msg.Picked)
	}

	return err

}
