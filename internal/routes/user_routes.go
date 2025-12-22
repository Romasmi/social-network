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
	profileFriendsRepo := repository.CreateProfileFriendsRepository(db)
	sessionRepo := repository.CreateSessionRepository(db)

	uow := repository.CreateUnitOfWork(db)

	userService := services.CreateUserService(cityRepo, profileRepo, profileFriendsRepo, uow)
	sessionService := services.CreateSessionService(sessionRepo)

	userHandler := user_handler.CreateUserHandler(userService)

	r.HandleFunc("/user/register", userHandler.RegisterUserHandler).Methods(http.MethodPost)

	privateRoute := r.PathPrefix("/").Subrouter()

	authMiddleware := middleware.CreateAuthMiddleware(sessionService)

	privateRoute.Use(authMiddleware.Process)
	privateRoute.HandleFunc("/user/get/{userId}", userHandler.GetUserHandler).Methods(http.MethodGet)
	privateRoute.HandleFunc("/user/search", userHandler.SearchUserHandler).Methods(http.MethodGet)
	privateRoute.HandleFunc("/friend/set/{userId}", userHandler.SetFriend).Methods(http.MethodPut)
	privateRoute.HandleFunc("/friend/delete/{userId}", userHandler.DeleteFriend).Methods(http.MethodPut)
}
