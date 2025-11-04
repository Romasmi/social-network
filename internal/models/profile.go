package models

type Gender string

const (
	Male   Gender = "male"
	FeMale Gender = "female"
)

type Profile struct {
	ID         string `json:"id"`
	UserId     string `json:"userId"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Birthdate  string `json:"birthdate"`
	Gender     Gender `json:"gender"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
}
