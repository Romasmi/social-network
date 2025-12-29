package services

import (
	"context"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/google/uuid"
)

type PostService struct {
	postsRepo repository.PostsRepository
}

func CreatePostService(postsRepo repository.PostsRepository) *PostService {
	return &PostService{postsRepo: postsRepo}
}

func (s *PostService) CreatePost(ctx context.Context, profileID uuid.UUID, text string) (*models.Post, error) {
	postID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	post := &models.Post{
		ID:        postID,
		ProfileId: profileID,
		Text:      text,
	}
	return s.postsRepo.CreatePost(ctx, post)
}

func (s *PostService) UpdatePost(ctx context.Context, profileID uuid.UUID, postID uuid.UUID, text string) (*models.Post, error) {
	return s.postsRepo.UpdatePost(ctx, postID, profileID, text)
}

func (s *PostService) DeletePost(ctx context.Context, profileID uuid.UUID, postID uuid.UUID) error {
	return s.postsRepo.DeletePost(ctx, postID, profileID)
}

func (s *PostService) GetPost(ctx context.Context, profileID uuid.UUID, postID uuid.UUID) (*models.Post, error) {
	return s.postsRepo.GetPost(ctx, postID, profileID)
}

func (s *PostService) GetFeed(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	return s.postsRepo.GetFeed(ctx, profileID, limit, offset)
}
