package app

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	router *mux.Router
}

func CreateApp() (*App, error) {
	app := &App{}
	err := app.init()
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) init() error {
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
