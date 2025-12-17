package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostsRepository struct {
	db DBQuerier
}

const postsTable = "posts"

func CreatePostsRepository(db DBQuerier) *PostsRepository {
	return &PostsRepository{db: db}
}

func (r *PostsRepository) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	const query = `
        INSERT INTO %s (id, profile_id, text)
        VALUES ($1, $2, $3)
        RETURNING id, profile_id, text
    `
	sql := fmt.Sprintf(query, postsTable)

	created := &models.Post{}
	err := r.db.QueryRow(ctx, sql, post.ID, post.ProfileId, post.Text).Scan(
		&created.ID,
		&created.ProfileId,
		&created.Text,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrDuplicate
			}
			return nil, err
		}
		return nil, err
	}
	return created, nil
}

func (r *PostsRepository) UpdatePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID, newText string) (*models.Post, error) {
	const query = `
        UPDATE %s
        SET text = $1
        WHERE id = $2 AND profile_id = $3
        RETURNING id, profile_id, text
    `
	sql := fmt.Sprintf(query, postsTable)

	updated := &models.Post{}
	err := r.db.QueryRow(ctx, sql, newText, postID, profileID).Scan(
		&updated.ID,
		&updated.ProfileId,
		&updated.Text,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return updated, nil
}

func (r *PostsRepository) DeletePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID) error {
	const query = `
        DELETE FROM %s
        WHERE id = $1 AND profile_id = $2
    `
	sql := fmt.Sprintf(query, postsTable)

	tag, err := r.db.Exec(ctx, sql, postID, profileID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostsRepository) GetPost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID) (*models.Post, error) {
	const query = `
        SELECT id, profile_id, text
        FROM %s
        WHERE id = $1 AND profile_id = $2
        LIMIT 1
    `
	sql := fmt.Sprintf(query, postsTable)

	post := &models.Post{}
	err := r.db.QueryRow(ctx, sql, postID, profileID).Scan(
		&post.ID,
		&post.ProfileId,
		&post.Text,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return post, nil
}
