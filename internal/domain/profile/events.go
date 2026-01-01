package profile

import (
	"github.com/Romasmi/social-network/internal/events"
	"github.com/google/uuid"
)

type FriendDeletedEventData struct {
	ProfileID uuid.UUID `json:"profileId"`
	FriendID  uuid.UUID `json:"friendId"`
}

func NewFriendDeletedEvent(data *FriendDeletedEventData) *events.Event {
	return events.NewEvent(events.FriendDeleted, data)
}
