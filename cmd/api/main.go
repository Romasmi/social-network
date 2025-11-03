package main

import (
	"fmt"
	"path/filepath"

	"github.com/Romasmi/social-network/internal/app"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	appInstance, err := app.CreateApp("../../")
	if err != nil {
		fmt.Printf("error while app init: %v", err)
	}
	defer appInstance.OnStop()

	path, _ := filepath.Abs("../../migrations")
	m, err := migrate.New(
		"file://"+path,
		appInstance.Config.Database.URL,
	)
	if m == nil || err != nil {
		fmt.Printf("unable to create migrations driver: %v\n", err)
		return
	}
	err = m.Up()
	if err != nil {
		fmt.Printf("error while running up migrations: %v\n", err)
		return
	}

	appInstance.Run()

	// TODO Add gracefull shutdown
}
