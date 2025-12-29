package commands

import (
	"github.com/Romasmi/social-network/internal/database"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/repository/posts_repository"
	"github.com/Romasmi/social-network/internal/services"
)

type ImportHandler struct {
	postService *services.PostService
	userService *services.UserService
}

type App interface {
	GetDB() *database.DbConnection
}

func CreateImportHandler(app App) *ImportHandler {
	postRepo := posts_repository.CreatePostsRepository(app.GetDB().DB)
	postService := services.CreatePostService(postRepo)

	cityRepo := repository.CreateCityRepository(app.GetDB().DB)
	profileRepo := repository.CreateProfileRepository(app.GetDB().DB)
	uow := repository.CreateUnitOfWork(app.GetDB().DB)
	userService := services.CreateUserService(cityRepo, profileRepo, nil, uow)

	return &ImportHandler{
		postService: postService,
		userService: userService,
	}
}
