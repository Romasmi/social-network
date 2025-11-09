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
	query := fmt.Sprintf(`
		INSERT INTO %s (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING *
	`, usersTable)

	newUser, err := scanUser(r.db.QueryRow(ctx, query, user.ID, user.Email, user.PasswordHash))
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
	return newUser, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := fmt.Sprintf(`
		SELECT *
        FROM %v
		WHERE email = $1
		LIMIT 1
	`, usersTable)

	user, err := scanUser(r.db.QueryRow(ctx, query, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to user url by %s: %w", email, err)
	}
	return user, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	query := fmt.Sprintf(`
		SELECT *
        FROM %v
		WHERE id = $1
		LIMIT 1
	`, usersTable)

	user, err := scanUser(r.db.QueryRow(ctx, query, userId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by %s: %w", userId, err)
	}
	return user, nil
}

func scanUser(row pgx.Row) (*models.User, error) {
	user := &models.User{}
	return user, row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.UpdateAt,
		&user.CreatedAt,
	)
}
