package models

import (
	"time"

	"github.com/google/uuid"
)

type City struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	CountryCode   string    `json:"countryCode"`
	StateProvince string    `json:"stateProvince"`
	CreatedAt     time.Time `json:"createdAt"`
}
