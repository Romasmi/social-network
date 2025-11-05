package user_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := uuid.MustParse(vars["userId"])

	profile, err := h.userService.GetUser(r.Context(), userId)
	// TODO check error if user not found
	if err != nil {
		fmt.Printf("error while user registration: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(profile)
	if err != nil {
		fmt.Printf("error while encoding response: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
