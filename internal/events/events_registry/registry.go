package events_registry

import (
	"context"

	"github.com/Romasmi/social-network/internal/events"
	eventhandler "github.com/Romasmi/social-network/internal/handlers/events"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/repository/posts_repository"
	"github.com/redis/go-redis/v9"
)

type EventHandler func(ctx context.Context, event *events.Event) error

type EventRegistry interface {
	Register(eventType events.EventType, handler EventHandler)
	GetHandlers(eventType events.EventType) []EventHandler
}

type EventRegistryImpl struct {
	handlers map[events.EventType][]EventHandler
}

func NewEventRegistry(db repository.DBQuerier, redis *redis.Client) EventRegistry {
	r := &EventRegistryImpl{
		handlers: make(map[events.EventType][]EventHandler),
	}
	r.registerEventHandlers(db, redis)
	return r
}

func (r *EventRegistryImpl) Register(eventType events.EventType, handler EventHandler) {
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

func (r *EventRegistryImpl) GetHandlers(eventType events.EventType) []EventHandler {
	return r.handlers[eventType]
}

func (r *EventRegistryImpl) registerEventHandlers(db repository.DBQuerier, redis *redis.Client) {
	profileFriendsRepo := repository.CreateProfileFriendsRepository(db)
	postsCache := posts_repository.NewPostsCache(redis)
	postsRepository := posts_repository.CreateCachedPostsRepository(posts_repository.CreatePostsRepository(db), redis)

	postHandler := eventhandler.CreatePostEventsHandler(profileFriendsRepo, postsCache)
	profileHandler := eventhandler.CreateProfileEventHandler(profileFriendsRepo, postsCache, postsRepository)

	r.Register(events.PostCreated, postHandler.OnPostCreated)
	r.Register(events.PostDeleted, postHandler.OnPostDeleted)

	r.Register(events.FriendDeleted, profileHandler.OnFriendDeleted)
}
