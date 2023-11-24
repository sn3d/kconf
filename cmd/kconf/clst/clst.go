package clst

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sn3d/kconf/cmd/kconf/clst/mod"
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/sn3d/kconf/pkg/tui/list"

	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "clst",
	Usage:     "change cluster for context",
	ArgsUsage: "[CLUSTER]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
		&cli.StringFlag{
			Name:    "context",
			Aliases: []string{"c"},
			Usage:   "context for which you want to change cluster",
		},
	},

	Subcommands: []*cli.Command{
		mod.Cmd,
	},

	// main entry point for 'export'
	Action: func(cCtx *cli.Context) error {
		cluster := cCtx.Args().First()
		if cluster != "" {
			return directChange(cCtx, cluster)
		} else {
			return showTUI(cCtx)
		}
	},
}

func directChange(cCtx *cli.Context, cluster string) error {
	kc, err := kconf.Open(cCtx.String("kubeconfig"))
	if err != nil {
		return err
	}

	err = kc.ChangeCluster(cCtx.String("context"), cluster)
	if err != nil {
		return err
	}

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

	model, err := NewModel(kc, cCtx.String("context"))
	if err != nil {
		return err
	}

	m, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	out := m.(Model)

	switch msg := out.ExitMsg.(type) {
	case list.SaveAndQuitMsg:
		fmt.Printf("changes saved\n")
	case list.PickedMsg:
		fmt.Printf("cluster changed to %s\n", msg.Picked)
	}

	return err
}
