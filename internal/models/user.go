package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	IsActive     bool      `json:"isActive"`
	PasswordHash string    `json:"-"`
	UpdateAt     time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

type CreateUserModel struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"` // format 2017-02-01
	Biography  string `json:"biography"`
	Gender     Gender `json:"gender"`
	City       string `json:"city"`
	Password   string `json:"password"`
}
