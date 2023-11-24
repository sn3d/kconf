package usr

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sn3d/kconf/cmd/kconf/usr/mod"
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/sn3d/kconf/pkg/tui/list"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "usr",
	Usage:     "change user for context",
	ArgsUsage: "[USER]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "path to kubeconfig from where export context",
		},
		&cli.StringFlag{
			Name:    "context",
			Aliases: []string{"c"},
			Usage:   "context for which you want to change user",
		},
	},

	Subcommands: []*cli.Command{
		mod.Cmd,
	},

	// main entry point
	Action: func(cCtx *cli.Context) error {
		user := cCtx.Args().First()
		if user != "" {
			return directChange(cCtx, user)
		} else {
			return showTUI(cCtx)
		}

	},
}

func directChange(cCtx *cli.Context, user string) error {

	kc, err := kconf.Open(cCtx.String("kubeconfig"))
	if err != nil {
		return err
	}

	err = kc.ChangeUser(cCtx.String("context"), user)
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

	out := m.(Model)
	switch msg := out.ExitMsg.(type) {
	case list.SaveAndQuitMsg:
		fmt.Printf("changes saved\n")
	case list.PickedMsg:
		fmt.Printf("user changed to %s\n", msg.Picked)
	}

	return err
}
