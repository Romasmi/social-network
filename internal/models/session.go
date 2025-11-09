package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID       `json:"id"`
	UserId    uuid.UUID       `json:"-"`
	Metadata  json.RawMessage `json:"-"`
	ExpiresAt time.Time       `json:"-"`
	CreatedAt time.Time       `json:"-"`
}
