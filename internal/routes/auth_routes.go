package routes

import (
	"net/http"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/handlers/auth_handler"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuthHandlers(r *mux.Router, db *pgxpool.Pool, cng *config.Config) {
	authHandler := auth_handler.CreateAuthHandler()
	r.HandleFunc("/login", authHandler.LoginHandler).Methods(http.MethodPost)
}
