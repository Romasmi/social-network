package posts_repository

import (
	"context"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type cachedPostsRepositoryImpl struct {
	postsRepo PostsRepository
	redis     *redis.Client
}

func CreateCachedPostsRepository(postsRepo PostsRepository, rds *redis.Client) PostsRepository {
	return &cachedPostsRepositoryImpl{postsRepo: postsRepo, redis: rds}
}

func (c cachedPostsRepositoryImpl) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	return c.postsRepo.CreatePost(ctx, post)
}

func (c cachedPostsRepositoryImpl) UpdatePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID, newText string) (*models.Post, error) {
	return c.postsRepo.UpdatePost(ctx, postID, profileID, newText)
}

func (c cachedPostsRepositoryImpl) DeletePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID) error {
	return c.postsRepo.DeletePost(ctx, postID, profileID)
}

func (c cachedPostsRepositoryImpl) GetPost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID) (*models.Post, error) {
	return c.postsRepo.GetPost(ctx, postID, profileID)
}

func (c cachedPostsRepositoryImpl) GetFeed(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	return c.postsRepo.GetFeed(ctx, profileID, limit, offset)
}
