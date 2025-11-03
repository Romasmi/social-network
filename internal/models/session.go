package models

type Session struct {
	ID       string      `json:"id"`
	UserId   string      `json:"userId"`
	Metadata interface{} `json:"metadata"`
}
