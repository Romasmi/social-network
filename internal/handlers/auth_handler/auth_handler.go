package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/Romasmi/social-network/internal/services"
	"github.com/Romasmi/social-network/internal/utils"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func CreateAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.ErrorInvalidRequestBody(w, err)
		return
	}

}
