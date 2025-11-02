package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/database"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	DbConn *database.DbConnection
	Config *config.Config
	router *mux.Router
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

	return nil
}

func (a *App) OnStop() {

}

func (a *App) Run() {
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	})
	headers := handlers.AllowedHeaders([]string{
		"Content-Type",
		"Authorization",
	})
	origins := handlers.AllowedOrigins([]string{"*"})

	// 		":"+strconv.Itoa(int(a.Config.Server.Port)),
	err := http.ListenAndServe(
		":8080",
		handlers.CORS(credentials, methods, origins, headers)(a.router))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
