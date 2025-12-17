package routes

import (
	"net/http"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/handlers/post_handler"
	"github.com/Romasmi/social-network/internal/middleware"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPostRoutes(r *mux.Router, db *pgxpool.Pool, cng *config.Config) {
	postsRepo := repository.CreatePostsRepository(db)
	sessionRepo := repository.CreateSessionRepository(db)

	postService := services.CreatePostService(postsRepo)
	sessionService := services.CreateSessionService(sessionRepo)

	postHandler := post_handler.CreatePostHandler(postService)

	privateRoute := r.PathPrefix("/").Subrouter()
	authMiddleware := middleware.CreateAuthMiddleware(sessionService)
	privateRoute.Use(authMiddleware.Process)

	privateRoute.HandleFunc("/post/create", postHandler.Create).Methods(http.MethodPost)
	privateRoute.HandleFunc("/post/update", postHandler.Update).Methods(http.MethodPut)
	privateRoute.HandleFunc("/post/delete/{postId}", postHandler.Delete).Methods(http.MethodDelete)
	privateRoute.HandleFunc("/post/get/{postId}", postHandler.Get).Methods(http.MethodGet)
}
