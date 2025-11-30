package cli

import (
	"context"
	"fmt"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/database"
	"github.com/spf13/cobra"
)

type App struct {
	DbConn *database.DbConnection
	Config *config.Config
	Cmd    *cobra.Command
}

func CreateApp(configPath string) (*App, error) {
	app := &App{}
	err := app.init(configPath)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) init(configPath string) error {
	envConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading Config: %v\n", err)
	}
	a.Config = envConfig

	dbConn := &database.DbConnection{Config: &envConfig.Database}
	if err = dbConn.Connect(); err != nil {
		return fmt.Errorf("error connecting to DB: %v\n", err)
	}

	a.DbConn = dbConn
	a.Cmd = &cobra.Command{
		Use:   "social-network",
		Short: "CLI for managing social network",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to CLI of Social network")
		},
	}
	a.iniCommands()

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var shutdownErr error

	if a.DbConn != nil && a.DbConn.DB != nil {
		fmt.Println("Closing database connections...")
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown timeout reached, forcing database close")
		default:
			a.DbConn.DB.Close()
		}
	}

	fmt.Println("Cleanup completed")
	return shutdownErr
}
