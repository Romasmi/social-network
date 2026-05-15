package posts_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/social-network/internal/domain/post"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, post *post.Post) (*post.Post, error)
	UpdatePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID, newText string) (*post.Post, error)
	DeletePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID) error
	GetPost(ctx context.Context, postID uuid.UUID) (*post.Post, error)
	GetPostsByIds(ctx context.Context, postID []uuid.UUID) ([]*post.Post, error)
	GetPostsByProfileId(ctx context.Context, profileID uuid.UUID) ([]*post.Post, error)
	GetFeed(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*post.Post, error)
}

type postsRepositoryImpl struct {
	writer repository.DBQuerier
	reader repository.DBQuerier
}

const postsTable = "posts"

func CreatePostsRepository(writer repository.DBQuerier, reader repository.DBQuerier) PostsRepository {
	return &postsRepositoryImpl{writer: writer, reader: reader}
}

func (r *postsRepositoryImpl) CreatePost(ctx context.Context, postModel *post.Post) (*post.Post, error) {
	const query = `
        INSERT INTO %s (id, profile_id, text)
        VALUES ($1, $2, $3)
        RETURNING id, profile_id, text
    `
	sql := fmt.Sprintf(query, postsTable)

	created := &post.Post{}
	err := r.writer.QueryRow(ctx, sql, postModel.ID, postModel.ProfileId, postModel.Text).Scan(
		&created.ID,
		&created.ProfileId,
		&created.Text,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, repository.ErrDuplicate
			}
			return nil, err
		}
		return nil, err
	}
	return created, nil
}

func (r *postsRepositoryImpl) UpdatePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID, newText string) (*post.Post, error) {
	const query = `
        UPDATE %s
        SET text = $1
        WHERE id = $2 AND profile_id = $3
        RETURNING id, profile_id, text
    `
	sql := fmt.Sprintf(query, postsTable)

	updated := &post.Post{}
	err := r.writer.QueryRow(ctx, sql, newText, postID, profileID).Scan(
		&updated.ID,
		&updated.ProfileId,
		&updated.Text,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return updated, nil
}

func (r *postsRepositoryImpl) DeletePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID) error {
	const query = `
        DELETE FROM %s
        WHERE id = $1 AND profile_id = $2
    `
	sql := fmt.Sprintf(query, postsTable)

	tag, err := r.writer.Exec(ctx, sql, postID, profileID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *postsRepositoryImpl) GetPost(ctx context.Context, postID uuid.UUID) (*post.Post, error) {
	const query = `
        SELECT id, profile_id, text
        FROM %s
        WHERE id = $1
        LIMIT 1
    `
	sql := fmt.Sprintf(query, postsTable)

	post := &post.Post{}
	err := r.reader.QueryRow(ctx, sql, postID).Scan(
		&post.ID,
		&post.ProfileId,
		&post.Text,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return post, nil
}

func (r *postsRepositoryImpl) GetPostsByProfileId(ctx context.Context, profileID uuid.UUID) ([]*post.Post, error) {
	const query = `
        SELECT id, profile_id, text
        FROM %s
        WHERE profile_id = $1
    `
	sql := fmt.Sprintf(query, postsTable)

	rows, err := r.reader.Query(ctx, sql, profileID)
	return r.rowsToPosts(rows, err)
}

func (r *postsRepositoryImpl) GetPostsByIds(ctx context.Context, postIDs []uuid.UUID) ([]*post.Post, error) {
	const query = `
        SELECT id, profile_id, text
        FROM %s
        WHERE id = ANY($1)
    `
	sql := fmt.Sprintf(query, postsTable)

	rows, err := r.reader.Query(ctx, sql, postIDs)
	return r.rowsToPosts(rows, err)
}

func (r *postsRepositoryImpl) GetFeed(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*post.Post, error) {
	const query = `
		SELECT *
		FROM %s
		WHERE profile_id IN (
			SELECT profile2_id
			FROM %s
			WHERE profile1_id = $1
			UNION ALL
			SELECT profile1_id
			FROM %s
			WHERE profile2_id = $1
		) AND profile_id != $1
		LIMIT $2
		OFFSET $3
    `
	sql := fmt.Sprintf(query, postsTable, repository.ProfileFriendsTable, repository.ProfileFriendsTable)

	rows, err := r.reader.Query(ctx, sql, profileID, limit, offset)
	return r.rowsToPosts(rows, err)
}

func (r *postsRepositoryImpl) rowsToPosts(rows pgx.Rows, err error) ([]*post.Post, error) {
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*post.Post{}, nil
		}
		return nil, err
	}
	posts := []*post.Post{}
	for rows.Next() {
		p := &post.Post{}
		err := rows.Scan(
			&p.ID,
			&p.ProfileId,
			&p.Text,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
