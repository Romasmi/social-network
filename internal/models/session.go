package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID       `json:"token"`
	UserId    uuid.UUID       `json:"-"`
	ProfileId uuid.UUID       `json:"-"`
	Metadata  json.RawMessage `json:"-"`
	ExpiresAt time.Time       `json:"expiresAt"`
	CreatedAt time.Time       `json:"-"`
}
