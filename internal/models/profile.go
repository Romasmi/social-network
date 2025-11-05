package models

import "github.com/google/uuid"

type Gender string

const (
	Male   Gender = "male"
	FeMale Gender = "female"
)

type Profile struct {
	ID         uuid.UUID `json:"id"`
	UserId     uuid.UUID `json:"userId"`
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"secondName"`
	Birthdate  string    `json:"birthdate"`
	Gender     Gender    `json:"gender"`
	Biography  string    `json:"biography"`
	City       string    `json:"city"`
	CityId     uuid.UUID `json:"-"`
}
