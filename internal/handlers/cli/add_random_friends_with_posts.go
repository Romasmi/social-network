package cli

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/Romasmi/social-network/internal/domain/profile"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (h *ImportHandler) AddRandomFriendsWithPosts(ctx context.Context, profileId uuid.UUID, numberOfFriends int) error {
	if numberOfFriends == 0 {
		return nil
	}

	candidates, err := h.userService.SearchUsers(ctx, &profile.UserSearchParams{})
	if err != nil {
		return err
	}
	ids := make([]uuid.UUID, 0, len(candidates))
	for _, p := range candidates {
		ids = append(ids, p.ID)
	}
	if len(ids) == 0 {
		return fmt.Errorf("no candidate profiles found")
	}

	rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	count := numberOfFriends
	if count > len(ids) {
		count = len(ids)
	}
	selected := ids[:count]

	var errors []string
	for _, fid := range selected {
		if err := h.userService.SetFriend(ctx, profileId, fid); err != nil {
			errors = append(errors, fmt.Sprintf("set friend %s: %v", fid, err))
			continue
		}
		if _, err := h.postService.CreatePost(ctx, fid, gofakeit.LoremIpsumSentence(100)); err != nil {
			errors = append(errors, fmt.Sprintf("create post for %s: %v", fid, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("completed with errors: %s", strings.Join(errors, ";\n"))
	}
	return nil
}
