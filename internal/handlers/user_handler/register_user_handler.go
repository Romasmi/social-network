package user_handler

import (
	"encoding/json"
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

func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.ErrorInvalidRequestBody(w)
		return
	}

}
