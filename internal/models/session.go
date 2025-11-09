package models

import "github.com/google/uuid"

type Session struct {
	ID       uuid.UUID   `json:"id"`
	UserId   string      `json:"-"`
	Metadata interface{} `json:"-"`
}
