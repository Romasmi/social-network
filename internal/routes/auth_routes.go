package routes

import (
	"net/http"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/handlers/auth_handler"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuthHandlers(r *mux.Router, db *pgxpool.Pool, cng *config.Config) {
	userRepo := repository.CreateUserRepository(db)
	sessionRepo := repository.CreateSessionRepository(db)
	authService := services.CreateAuthService(userRepo, sessionRepo)

	authHandler := auth_handler.CreateAuthHandler(authService)
	r.HandleFunc("/login", authHandler.LoginHandler).Methods(http.MethodPost)
}
