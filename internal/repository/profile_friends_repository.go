package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProfileFriendsRepository struct {
	db DBQuerier
}

const ProfileFriendsTable = "profile_friends"

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
	sql := fmt.Sprintf(query, ProfileFriendsTable)

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
	sql := fmt.Sprintf(query, ProfileFriendsTable)

	_, err := r.db.Exec(
		ctx,
		sql,
		profile1Id,
		profile2Id,
	)
	return err
}

func (r *ProfileFriendsRepository) GetFriendsIds(ctx context.Context, profileID uuid.UUID) ([]uuid.UUID, error) {
	const query = `
			SELECT profile2_id
			FROM %s
			WHERE profile1_id = $1
			UNION
			SELECT profile1_id
			FROM %s
			WHERE profile2_id = $1
	`
	sql := fmt.Sprintf(query, ProfileFriendsTable)
	rows, err := r.db.Query(ctx, sql, profileID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []uuid.UUID{}, nil
		}
		return nil, err
	}
	ids := []uuid.UUID{}
	for rows.Next() {
		id := uuid.UUID{}
		err := rows.Scan(
			&id,
		)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
