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

func (r *ProfileFriendsRepository) SetFriend(ctx context.Context, profile1Id, profile2Id uuid.UUID) error {
	const query = `
		INSERT INTO %s (profile1_id, profile2_id)
		VALUES ($1, $2)
		ON CONFLICT (profile1_id, profile2_id) DO NOTHING
	`
	sql := fmt.Sprintf(query, profileFriendsTable)

	_, err := r.db.Exec(
		ctx,
		sql,
		profile1Id,
		profile2Id,
	)
	return err
}

func (r *ProfileFriendsRepository) DeleteFriend(ctx context.Context, profile1Id, profile2Id uuid.UUID) error {
	const query = `
		DELETE 
		FROM %s
		WHERE profile1_id = $1 AND profile2_id = $2 OR profile2_id = $1 AND profile1_id = $2
	`
	sql := fmt.Sprintf(query, profileFriendsTable)

	_, err := r.db.Exec(
		ctx,
		sql,
		profile1Id,
		profile2Id,
	)
	return err
}
