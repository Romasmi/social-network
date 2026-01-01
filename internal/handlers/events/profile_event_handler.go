package events

import (
	"context"
	"fmt"

	"github.com/Romasmi/social-network/internal/domain/profile"
	"github.com/Romasmi/social-network/internal/events"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/repository/posts_repository"
	"github.com/google/uuid"
)

type ProfileEventHandler struct {
	profileFriendRepo *repository.ProfileFriendsRepository
	postsRepo         posts_repository.PostsRepository
	postsCache        *posts_repository.PostsCache
}

func CreateProfileEventHandler(
	profileFriendRepo *repository.ProfileFriendsRepository,
	postsCache *posts_repository.PostsCache,
	postsRepo posts_repository.PostsRepository,
) *ProfileEventHandler {
	return &ProfileEventHandler{
		profileFriendRepo: profileFriendRepo,
		postsCache:        postsCache,
		postsRepo:         postsRepo,
	}
}

func (h *ProfileEventHandler) OnFriendAdded(ctx context.Context, event *events.Event) error {
	data, err := events.CastDataTo[profile.FriendDeletedEventData](event)
	if err != nil {
		return fmt.Errorf("invalid event data: %v", event)
	}
	err = h.pushUserPostsFromOthersFeed(ctx, data.ProfileID, data.FriendID)
	if err != nil {
		return err
	}
	return h.pushUserPostsFromOthersFeed(ctx, data.FriendID, data.ProfileID)
}

func (h *ProfileEventHandler) OnFriendDeleted(ctx context.Context, event *events.Event) error {
	data, err := events.CastDataTo[profile.FriendDeletedEventData](event)
	if err != nil {
		return fmt.Errorf("invalid event data: %v", event)
	}
	err = h.deleteUserPostsFromOthersFeed(ctx, data.ProfileID, data.FriendID)
	if err != nil {
		return err
	}
	return h.deleteUserPostsFromOthersFeed(ctx, data.FriendID, data.ProfileID)
}

func (h *ProfileEventHandler) deleteUserPostsFromOthersFeed(ctx context.Context, profileID, otherID uuid.UUID) error {
	posts, err := h.postsRepo.GetPostsByProfileId(ctx, profileID)
	if err != nil {
		return fmt.Errorf("getting posts by profile id: %w", err)
	}
	postIds := make([]uuid.UUID, len(posts))
	for i, v := range posts {
		postIds[i] = v.ID
	}
	return h.postsCache.DeleteFromFeeds(ctx, []uuid.UUID{otherID}, postIds)
}

func (h *ProfileEventHandler) pushUserPostsFromOthersFeed(ctx context.Context, profileID, otherID uuid.UUID) error {
	posts, err := h.postsRepo.GetPostsByProfileId(ctx, profileID)
	if err != nil {
		return fmt.Errorf("getting posts by profile id: %w", err)
	}
	return h.postsCache.PushToFeeds(ctx, []uuid.UUID{otherID}, posts)
}
