package post

import (
	"github.com/Romasmi/social-network/internal/events"
	"github.com/google/uuid"
)

type (
	DeletedEventData struct {
		ProfileID uuid.UUID `json:"profileId"`
		PostID    uuid.UUID `json:"postId"`
	}

	CreatedEventData struct {
		Post *Post `json:"post"`
	}
)

func NewPostCreatedEvent(data *CreatedEventData) *events.Event {
	return events.NewEvent(events.PostCreated, data)
}

func NewPostDeletedEvent(data *DeletedEventData) *events.Event {
	return events.NewEvent(events.PostDeleted, data)
}
