package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ProfileFriendsRepository struct {
	db DBQuerier
}

const profileFriendsTable = "profile_friends"

func CreateProfileFriendsRepository(db DBQuerier) *ProfileFriendsRepository {
	return &ProfileFriendsRepository{
		db: db,
	}
}

func (r *ProfileFriendsRepository) AddFriend(ctx context.Context, profile1Id, profile2Id uuid.UUID) error {
	const query = `
		INSERT INTO %s (profile1_id, profile2_id)
		VALUES ($1, $2)
	`
	sql := fmt.Sprintf(query, profileFriendsTable)

	_, err := r.db.Query(
		ctx,
		sql,
		profile1Id,
		profile2Id,
	)
	return err
}
