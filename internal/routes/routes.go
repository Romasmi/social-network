package routes

import (
	"fmt"
	"net/http"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/middleware"
	"github.com/Romasmi/social-network/internal/utils"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotFoundResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(
	router *mux.Router,
	db *pgxpool.Pool,
	config *config.Config,
) {
	if router == nil {
		panic("router must be initialized before routes registration")
	}
	router.Use(middleware.ResponseHeadersMiddleware)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	RegisterAuthHandlers(router, db, config)
	RegisterUserRoutes(router, db, config)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.JsonError(w, http.StatusNotFound, fmt.Errorf("undefined route"))
}
