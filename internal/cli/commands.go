package cli

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func (a *App) Execute() {
	if err := a.Cmd.Execute(); err != nil {
		fmt.Println(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func (a *App) iniCommands() {
	a.Cmd.AddCommand(&cobra.Command{
		Use:   "import",
		Short: "Import users",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] != "l" {
				return fmt.Errorf("Invalid argument: %v\n", args[0])
			}
			_, err := url.ParseRequestURI(args[1])
			if err != nil {
				return err
			}

			return a.importUserByLink(context.Background(), args[1])
		},
	})
}
