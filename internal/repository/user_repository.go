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
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

const usersTable = "users"

func CreateUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	const query = `
		INSERT INTO %s (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, created_at
	`
	sql := fmt.Sprintf(query, usersTable)

	newUser := &models.User{}
	err := r.db.QueryRow(ctx, sql, user.ID, user.Email, user.PasswordHash).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrDuplicate
			}
			return nil, err
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return newUser, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at
        FROM %s
		WHERE id = $1
	`
	sql := fmt.Sprintf(query, usersTable)

	user := &models.User{}
	err := r.db.QueryRow(ctx, sql, userId).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetUserByProfileId(ctx context.Context, profileId uuid.UUID) (*models.User, error) {
	const query = `
		SELECT id, first_name, last_name, city, created_at
        FROM %s
		WHERE id = (SELECT user_id FROM %s WHERE id = $1)
	`
	sql := fmt.Sprintf(query, usersTable, profilesTable)
	user := &models.User{}
	err := r.db.QueryRow(ctx, sql, profileId).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
