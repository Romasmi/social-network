package events

import (
	"context"
	"fmt"

	"github.com/Romasmi/social-network/internal/domain/post"
	"github.com/Romasmi/social-network/internal/events"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/repository/posts_repository"
	"github.com/google/uuid"
)

type PostEventHandler struct {
	profileFriendRepo *repository.ProfileFriendsRepository
	postsCache        *posts_repository.PostsCache
}

func CreatePostEventsHandler(
	profileFriendRepo *repository.ProfileFriendsRepository,
	postsCache *posts_repository.PostsCache,
) *PostEventHandler {
	return &PostEventHandler{profileFriendRepo: profileFriendRepo, postsCache: postsCache}
}

func (h *PostEventHandler) OnPostCreated(ctx context.Context, event *events.Event) error {
	data, err := events.CastDataTo[post.CreatedEventData](event)
	if err != nil {
		return fmt.Errorf("invalid event data: %v", event)
	}
	friendIds, err := h.profileFriendRepo.GetFriendsIds(ctx, data.Post.ProfileId)
	if err != nil {
		return fmt.Errorf("getting friends ids: %w", err)
	}
	return h.postsCache.PushToFeeds(ctx, friendIds, []*post.Post{data.Post})
}

func (h *PostEventHandler) OnPostDeleted(ctx context.Context, event *events.Event) error {
	data, err := events.CastDataTo[post.DeletedEventData](event)
	if err != nil {
		return fmt.Errorf("invalid event data: %v", event)
	}
	friendIds, err := h.profileFriendRepo.GetFriendsIds(ctx, data.ProfileID)
	if err != nil {
		return fmt.Errorf("getting friends ids: %w", err)
	}
	return h.postsCache.DeleteFromFeeds(ctx, friendIds, []uuid.UUID{data.PostID})
}
