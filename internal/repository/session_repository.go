package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

const sessionsTable = "sessions"

func CreateSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(ctx context.Context, s *models.Session) (*models.Session, error) {
	const query = `
		INSERT INTO %s (id, user_id, metadata, expires_at)
		VALUES ($1, $2, $3::jsonb, $4)
		RETURNING id, user_id, metadata, expires_at, created_at
	`
	sql := fmt.Sprintf(query, sessionsTable)

	created := &models.Session{}
	err := r.db.QueryRow(ctx, sql, s.ID, s.UserId, s.Metadata, s.ExpiresAt).Scan(
		&created.ID,
		&created.UserId,
		&created.Metadata,
		&created.ExpiresAt,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return created, nil
}

func (r *SessionRepository) GetSessionById(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	const query = `
		SELECT s.id, s.user_id, s.metadata, s.expires_at, s.created_at, p.id as profile_id
		FROM %s AS s
			LEFT JOIN profiles AS p ON p.user_id = s.user_id 
		WHERE s.id = $1
		LIMIT 1
	`
	sql := fmt.Sprintf(query, sessionsTable)

	s := &models.Session{}
	err := r.db.QueryRow(ctx, sql, id).Scan(
		&s.ID,
		&s.UserId,
		&s.Metadata,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.ProfileId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get session by id: %w", err)
	}
	return s, nil
}
