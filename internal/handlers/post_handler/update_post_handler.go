package post_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/utils"
	"github.com/google/uuid"
)

type UpdatePostRequest struct {
	PostID string `json:"postId"`
	Text   string `json:"text"`
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	profileId, _ := uuid.Parse(r.Context().Value("profileId").(string))

	var payload UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.ErrorInvalidRequestBody(w, err)
		return
	}
	if payload.Text == "" || payload.PostID == "" {
		utils.ErrorInvalidRequestBody(w, fmt.Errorf("postId and text are required"))
		return
	}
	postId, err := uuid.Parse(payload.PostID)
	if err != nil {
		utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid postId"))
		return
	}

	post, err := h.postService.UpdatePost(r.Context(), profileId, postId, payload.Text)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.JsonErrorNotFound(w)
			return
		}
		fmt.Printf("error while updating post: %v\n", err)
		utils.JsonInternalServerError(w)
		return
	}
	utils.SuccessJsonResponse(w, post)
}
