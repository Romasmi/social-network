package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	IsActive     bool      `json:"isActive" db:"is_active"`
	PasswordHash string    `json:"-" db:"password_hash"`
	UpdateAt     time.Time `json:"updatedAt" db:"updated_at"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
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
