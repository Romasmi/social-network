package routes

import (
	"net/http"

	"github.com/Romasmi/social-network/internal/events/publisher"
	"github.com/Romasmi/social-network/internal/handlers/http/user_handler"
	"github.com/Romasmi/social-network/internal/middleware"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/services"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router, writer, reader repository.DBQuerier, publisher publisher.Publisher) {
	cityRepo := repository.CreateCityRepository(writer, reader)
	profileRepo := repository.CreateProfileRepository(writer, reader)
	profileFriendsRepo := repository.CreateProfileFriendsRepository(writer, reader)
	sessionRepo := repository.CreateSessionRepository(writer, reader)

	uow := repository.CreateUnitOfWork(writer, reader)

	userService := services.CreateUserService(cityRepo, profileRepo, profileFriendsRepo, uow, publisher)
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
