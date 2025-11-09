package models

import (
	"time"

	"github.com/google/uuid"
)

type City struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	CountryCode   string    `json:"countryCode" db:"country_code"`
	StateProvince string    `json:"stateProvince" db:"state_province"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}
