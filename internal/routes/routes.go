package routes

import (
	"fmt"
	"net/http"

	"github.com/Romasmi/social-network/internal/events/publisher"
	"github.com/Romasmi/social-network/internal/infra/database"
	"github.com/Romasmi/social-network/internal/infra/redis_client"
	"github.com/Romasmi/social-network/internal/metrics"
	"github.com/Romasmi/social-network/internal/middleware"
	"github.com/Romasmi/social-network/internal/utils"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App interface {
	GetWriter() *database.Connection
	GetRedis() *redis_client.Connection
	GetPublisher() publisher.Publisher
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
	router.Use(metrics.Middleware)
	router.Use(middleware.LoggingMiddleware)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	writer := app.GetWriter().Writer()
	reader := app.GetWriter().Reader()

	RegisterAuthHandlers(router, writer, reader)
	RegisterUserRoutes(router, writer, reader, app.GetPublisher())
	RegisterPostRoutes(router, writer, reader, app.GetRedis().Rdb, app.GetPublisher())
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	utils.JsonError(w, http.StatusNotFound, fmt.Errorf("undefined route"))
}
