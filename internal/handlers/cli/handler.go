package cli

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Romasmi/social-network/internal/events/publisher"
	"github.com/Romasmi/social-network/internal/infra/database"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/repository/posts_repository"
	"github.com/Romasmi/social-network/internal/services"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type ImportHandler struct {
	postService *services.PostService
	userService *services.UserService
}

type App interface {
	GetWriter() *database.Connection
	GetPublisher() publisher.Publisher
}

func RegisterCommands(cmd *cobra.Command, app App) {
	handler := CreateImportHandler(app)

	cmd.AddCommand(&cobra.Command{
		Use:   "import",
		Short: "Import users",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := url.ParseRequestURI(args[0])
			if err != nil {
				return err
			}
			start := time.Now()
			maxCount, _ := strconv.Atoi(args[1])
			err = handler.ImportUserByLink(context.Background(), args[0], maxCount)
			fmt.Printf("Import took %v\n", time.Since(start))
			return err
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:     "importPosts",
		Short:   "Import posts",
		Long:    "Import posts to user. If multiple users provide then posts will be assigned randomly to provided users",
		Example: "ImportPosts https://posts.com/json 019ad582-424c-7316-b98c-eee5bf0f8d08 019ad582-42b3-7340-a5b8-6025ce042c01",
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			_, err := url.ParseRequestURI(args[0])
			if err != nil {
				return err
			}
			start := time.Now()
			err = handler.ImportPosts(context.Background(), args[0], args[1:])
			fmt.Printf("Import took %v\n", time.Since(start))
			return err
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "addRandomFriendsWithPosts <profileId> <numberOfFriends>",
		Short:   "Add random friends to a profile and create a post for each",
		Example: "AddRandomFriendsWithPosts 019ad582-424c-7316-b98c-eee5bf0f8d08 10",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileID, err := uuid.Parse(args[0])
			if err != nil {
				return fmt.Errorf("invalid profileId: %w", err)
			}
			n, err := strconv.Atoi(args[1])
			if err != nil || n < 0 {
				return fmt.Errorf("invalid numberOfFriends: %v", args[1])
			}
			start := time.Now()
			err = handler.AddRandomFriendsWithPosts(context.Background(), profileID, n)
			fmt.Printf("Command took %v\n", time.Since(start))
			return err
		},
	})
}

func CreateImportHandler(app App) *ImportHandler {
	writer := app.GetWriter().Writer()
	reader := app.GetWriter().Reader()

	postRepo := posts_repository.CreatePostsRepository(writer, reader)
	postService := services.CreatePostService(postRepo, app.GetPublisher())

	cityRepo := repository.CreateCityRepository(writer, reader)
	profileRepo := repository.CreateProfileRepository(writer, reader)
	profileFriendsRepo := repository.CreateProfileFriendsRepository(writer, reader)
	uow := repository.CreateUnitOfWork(writer, reader)
	userService := services.CreateUserService(cityRepo, profileRepo, profileFriendsRepo, uow, app.GetPublisher())

	return &ImportHandler{
		postService: postService,
		userService: userService,
	}
}
