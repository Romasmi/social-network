package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Romasmi/social-network/internal/app"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	appInstance, err := app.CreateApp("../../")
	if err != nil {
		fmt.Printf("error while app init: %v", err)
		os.Exit(1)
	}

	path, _ := filepath.Abs("../../migrations")
	m, err := migrate.New(
		"file://"+path,
		appInstance.Config.Database.URL,
	)
	if m == nil || err != nil {
		fmt.Printf("unable to create migrations driver: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
			fmt.Printf("Error closing migration driver - source: %v, db: %v\n", sourceErr, dbErr)
		}
	}()

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Printf("error while running up migrations: %v\n", err)
		os.Exit(1)
	}

	go func() {
		if err := appInstance.Run(); err != nil {
			fmt.Printf("server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := appInstance.Shutdown(ctx); err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	}

	fmt.Println("Server stopped")
}
