package models

import (
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	Male   Gender = "male"
	FeMale Gender = "female"
)

type Profile struct {
	ID         uuid.UUID `json:"id"`
	UserId     uuid.UUID `json:"-"`
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"secondName"`
	Birthdate  time.Time `json:"birthdate"`
	Biography  string    `json:"biography"`
	Gender     Gender    `json:"gender"`
	City       string    `json:"city"`
	CityId     uuid.UUID `json:"-"`
}

type CreateProfileModel struct {
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	Birthdate  time.Time `json:"birthdate"` // format 2017-02-01
	Biography  string    `json:"biography"`
	Gender     Gender    `json:"gender"`
	City       string    `json:"city"`
	Password   string    `json:"password"`
}

type UserSearchParams struct {
	FirstName  string
	SecondName string
}
