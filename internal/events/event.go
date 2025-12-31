package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventType string
type EventData map[string]any

const (
	PostCreated   EventType = "post.created"
	PostDeleted   EventType = "post.deleted"
	FriendDeleted EventType = "friend.deleted"
)

type Event struct {
	ID        string            `json:"id"`
	Type      EventType         `json:"type"`
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Version   string            `json:"version"`
	Data      EventData         `json:"data"`
	Metadata  map[string]string `json:"metadata"`
}

func NewEvent(eventType EventType, data EventData) *Event {
	id, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
	}

	return &Event{
		ID:        id.String(),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Source:    "social-network-api",
		Version:   "1.0",
		Data:      data,
		Metadata:  make(map[string]string),
	}
}

func (e *Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func FromJSON(data []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
