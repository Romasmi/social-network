package app

import (
	"context"
	"fmt"
	"os"

	"github.com/Romasmi/social-network/internal/handlers/cli"
	"github.com/spf13/cobra"
)

type Cli struct {
	App *App
	Cmd *cobra.Command
}

func NewCli(app *App) *Cli {
	cli := &Cli{App: app}
	cli.init()
	return cli
}

func (a *Cli) init() {
	a.Cmd = &cobra.Command{
		Use:   "social-network",
		Short: "CLI for managing social network",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to CLI of Social network")
		},
	}
	cli.RegisterCommands(a.Cmd, a.App)
}

func (a *Cli) Run() {
	if err := a.Cmd.Execute(); err != nil {
		fmt.Println(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	if err := a.Shutdown(context.Background()); err != nil {
		fmt.Println(os.Stderr, "Shutdown error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func (a *Cli) Shutdown(ctx context.Context) error {
	return a.App.Shutdown(ctx)
}
