package user_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/Romasmi/social-network/internal/utils"
)

type CreateUserRequest struct {
	FirstName  string        `json:"first_name"`
	SecondName string        `json:"second_name"`
	Birthdate  string        `json:"birthdate"` // format 2017-02-01
	Biography  string        `json:"biography"`
	Gender     models.Gender `json:"gender"`
	City       string        `json:"city"`
	Password   string        `json:"password"`
}

func (payload *CreateUserRequest) toModel() *models.CreateUserModel {
	return &models.CreateUserModel{
		FirstName:  payload.FirstName,
		SecondName: payload.SecondName,
		Birthdate:  payload.Birthdate,
		Biography:  payload.Biography,
		Gender:     payload.Gender,
		City:       payload.City,
		Password:   payload.Password,
	}
}

func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.ErrorInvalidRequestBody(w)
		return
	}
	profile, err := h.userService.RegisterUser(r.Context(), payload.toModel())
	if err != nil {
		fmt.Printf("error while user registration: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(profile)
	if err != nil {
		fmt.Printf("error while encoding response: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
