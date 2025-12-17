package post_handler

import "github.com/Romasmi/social-network/internal/services"

type PostHandler struct {
	postService *services.PostService
}

func CreatePostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}
