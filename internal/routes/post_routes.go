package routes

import (
	"net/http"

	"github.com/Romasmi/social-network/internal/events"
	"github.com/Romasmi/social-network/internal/handlers/http/post_handler"
	"github.com/Romasmi/social-network/internal/middleware"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/repository/posts_repository"
	"github.com/Romasmi/social-network/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func RegisterPostRoutes(r *mux.Router, db *pgxpool.Pool, rds *redis.Client, publisher events.Publisher) {
	postsRepo := posts_repository.CreateCachedPostsRepository(
		posts_repository.CreatePostsRepository(db),
		rds,
		publisher,
	)
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
	privateRoute.HandleFunc("/post/feed", postHandler.GetFeed).Methods(http.MethodGet)
}
