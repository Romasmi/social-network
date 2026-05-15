package routes

import (
	"net/http"

	"github.com/Romasmi/social-network/internal/handlers/http/auth_handler"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/services"
	"github.com/gorilla/mux"
)

func RegisterAuthHandlers(r *mux.Router, writer, reader repository.DBQuerier) {
	userRepo := repository.CreateUserRepository(writer, reader)
	sessionRepo := repository.CreateSessionRepository(writer, reader)
	authService := services.CreateAuthService(userRepo, sessionRepo)

	authHandler := auth_handler.CreateAuthHandler(authService)
	r.HandleFunc("/login", authHandler.LoginHandler).Methods(http.MethodPost)
}
