package user_handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type userSearchParameters struct {
	firstName  string
	secondName string
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["userId"])
	if err != nil {
		utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid user id"))
		return
	}

	profile, err := h.userService.GetUserByProfileId(r.Context(), userId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.JsonErrorNotFound(w)
			return
		}

		fmt.Printf("error while retreiving a user: %v\n", err)
		utils.JsonInternalServerError(w)
		return
	}

	utils.SuccessJsonResponse(w, profile)
}

func (h *UserHandler) SearchUserHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userSearchData := &models.UserSearchParams{
		FirstName:  queryParams.Get("first_name"),
		SecondName: queryParams.Get("last_name"),
	}

	profiles, err := h.userService.SearchUsers(r.Context(), userSearchData)
	if err != nil {
		fmt.Printf("error while searching users: %v\n", err)
		utils.JsonInternalServerError(w)
		return
	}

	utils.SuccessJsonResponse(w, profiles)
}
