package events_registry

import (
	"context"

	"github.com/Romasmi/social-network/internal/app"
	"github.com/Romasmi/social-network/internal/events"
	eventhandler "github.com/Romasmi/social-network/internal/handlers/events"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/repository/posts_repository"
)

type EventHandler func(ctx context.Context, event *events.Event) error

type EventRegistry interface {
	Register(eventType events.EventType, handler EventHandler)
	GetHandlers(eventType events.EventType) []EventHandler
}

type EventRegistryImpl struct {
	handlers map[events.EventType][]EventHandler
}

func NewEventRegistry() EventRegistry {
	return &EventRegistryImpl{
		handlers: make(map[events.EventType][]EventHandler),
	}
}

func (r *EventRegistryImpl) Register(eventType events.EventType, handler EventHandler) {
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

func (r *EventRegistryImpl) GetHandlers(eventType events.EventType) []EventHandler {
	return r.handlers[eventType]
}

func RegisterEventHandlers(registry EventRegistry, app app.App) {
	postHandler := eventhandler.CreatePostEventsHandler(
		repository.CreateProfileFriendsRepository(app.GetDB().DB),
		posts_repository.NewPostsCache(app.GetRedis().Rdb),
	)
	registry.Register(events.PostCreated, postHandler.OnPostCreated)
	registry.Register(events.PostDeleted, postHandler.OnPostDeleted)
}
