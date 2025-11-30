package user_handler

import (
	"errors"
	"fmt"
	"net/http"

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
	userSearchData := &userSearchParameters{
		firstName:  queryParams.Get("first_name"),
		secondName: queryParams.Get("second_name"),
	}
	fmt.Println("search data", userSearchData)
}
