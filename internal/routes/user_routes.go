package routes

import (
	"net/http"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/handlers/user_handler"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(router *mux.Router, db *pgxpool.Pool, cng *config.Config) {
	cityRepo := repository.CreateCityRepository(db)
	uow := repository.CreateUnitOfWork(db)
	userService := services.CreateUserService(cityRepo, uow)
	userHandler := user_handler.CreateUserHandler(userService)

	privateRouter := router

	privateRouter.HandleFunc("/user/register", userHandler.RegisterUserHandler).Methods(http.MethodPost)
	privateRouter.HandleFunc("/user/get/{userId}", userHandler.GetUserHandler).Methods(http.MethodGet)
}
