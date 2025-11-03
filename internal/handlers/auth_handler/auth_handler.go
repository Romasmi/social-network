package auth_handler

import "net/http"

type AuthHandler struct {
}

func CreateAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

}
