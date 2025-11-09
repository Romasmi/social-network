package routes

import (
	"net/http"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/handlers/user_handler"
	"github.com/Romasmi/social-network/internal/middleware"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool, cng *config.Config) {
	cityRepo := repository.CreateCityRepository(db)
	profileRepo := repository.CreateProfileRepository(db)
	uow := repository.CreateUnitOfWork(db)
	userService := services.CreateUserService(cityRepo, profileRepo, uow)
	userHandler := user_handler.CreateUserHandler(userService)

	r.HandleFunc("/user/register", userHandler.RegisterUserHandler).Methods(http.MethodPost)

	privateRoute := r.PathPrefix("/").Subrouter()
	privateRoute.Use(middleware.AuthMiddleware)
	privateRoute.HandleFunc("/user/get/{userId}", userHandler.GetUserHandler).Methods(http.MethodGet)
}
