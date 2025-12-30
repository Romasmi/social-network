package post_handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	profileId, _ := uuid.Parse(r.Context().Value("profileId").(string))

	vars := mux.Vars(r)
	postIdStr := vars["postId"]
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid post id"))
		return
	}

	if err := h.postService.DeletePost(r.Context(), profileId, postId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.JsonErrorNotFound(w)
			return
		}
		fmt.Printf("error while deleting post: %v\n", err)
		utils.JsonInternalServerError(w)
		return
	}

	utils.SuccessJsonResponse(w, struct{}{})
}
