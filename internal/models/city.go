package models

import "github.com/google/uuid"

type City struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	CountyCode    string    `json:"countyCode"`
	StateProvince string    `json:"stateProvince"`
}
