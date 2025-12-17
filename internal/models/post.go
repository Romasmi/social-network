package models

import "github.com/google/uuid"

type Post struct {
	ID        uuid.UUID `json:"id"`
	ProfileId uuid.UUID `json:"profileId"`
	Text      string    `json:"text"`
}
