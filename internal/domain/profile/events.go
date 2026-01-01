package profile

import (
	"github.com/Romasmi/social-network/internal/events"
	"github.com/google/uuid"
)

type FriendAddedEventData struct {
	ProfileID uuid.UUID `json:"profileId"`
	FriendID  uuid.UUID `json:"friendId"`
}

type FriendDeletedEventData struct {
	ProfileID uuid.UUID `json:"profileId"`
	FriendID  uuid.UUID `json:"friendId"`
}

func NewFriendAddedEvent(data *FriendAddedEventData) *events.Event {
	return events.NewEvent(events.FriendAdded, data)
}

func NewFriendDeletedEvent(data *FriendDeletedEventData) *events.Event {
	return events.NewEvent(events.FriendDeleted, data)
}
