package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/database"
	"github.com/Romasmi/social-network/internal/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	DbConn *database.DbConnection
	Config *config.Config
	router *mux.Router
	server *http.Server
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

	router := mux.NewRouter()

	routes.RegisterRoutes(router, dbConn.DB, envConfig)

	a.router = router
	return nil
}

func (a *App) OnStop() {

}

func (a *App) Run() error {
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

	a.server = &http.Server{
		Addr:    ":" + strconv.Itoa(int(a.Config.Server.Port)),
		Handler: handlers.CORS(credentials, methods, origins, headers)(a.router),
	}

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var err error

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			err = fmt.Errorf("server shutdown error: %w", err)
		}
	}

	if a.DbConn != nil && a.DbConn.DB != nil {
		a.DbConn.DB.Close()
	}

	return err
}
