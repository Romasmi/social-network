package routes

import (
	"fmt"
	"net/http"

	"github.com/Romasmi/social-network/internal/infra/database"
	"github.com/Romasmi/social-network/internal/infra/redis"
	"github.com/Romasmi/social-network/internal/middleware"
	"github.com/Romasmi/social-network/internal/utils"
	"github.com/gorilla/mux"
)

type App interface {
	GetDB() *database.DbConnection
	GetRedis() *redis.Connection
}

type NotFoundResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(
	router *mux.Router,
	app App,
) {
	if router == nil {
		panic("router must be initialized before routes registration")
	}
	router.Use(middleware.ResponseHeadersMiddleware)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	RegisterAuthHandlers(router, app.GetDB().DB)
	RegisterUserRoutes(router, app.GetDB().DB)
	RegisterPostRoutes(router, app.GetDB().DB, app.GetRedis().Rdb)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	utils.JsonError(w, http.StatusNotFound, fmt.Errorf("undefined route"))
}
