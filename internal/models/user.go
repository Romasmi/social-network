package models

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	IsActive bool   `json:"isActive"`
}
