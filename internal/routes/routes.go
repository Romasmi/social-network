package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Romasmi/social-network/internal/config"
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

	RegisterUserRoutes(router, db, config)

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	response := &NotFoundResponse{Error: "Route Not found"}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
