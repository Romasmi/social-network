package post_handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIdStr := vars["postId"]
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid post id"))
		return
	}

	post, err := h.postService.GetPost(r.Context(), postId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.JsonErrorNotFound(w)
			return
		}
		fmt.Printf("error while getting post: %v\n", err)
		utils.JsonInternalServerError(w)
		return
	}

	utils.SuccessJsonResponse(w, post)
}

func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	profileId := uuid.MustParse(r.Context().Value("profileId").(string))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	posts, err := h.postService.GetFeed(r.Context(), profileId, limit, offset)
	if err != nil {
		fmt.Printf("error while getting feed: %v\n", err)
		utils.JsonInternalServerError(w)
		return
	}
	utils.SuccessJsonResponse(w, posts)
}
