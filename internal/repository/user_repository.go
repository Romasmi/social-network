package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/social-network/internal/models"
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
		INSERT INTO %s (id, email, password_hash, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`, usersTable)

	var newUser *models.User

	err := r.db.QueryRow(ctx, query, user.ID, user.Email, user.PasswordHash).Scan(&newUser)
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
		WHERE email = '$1'
		LIMIT 1
	`, usersTable)

	var user *models.User

	err := r.db.QueryRow(ctx, query, email).Scan(&user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get url by %s: %w", email, err)
	}
	return user, nil
}
