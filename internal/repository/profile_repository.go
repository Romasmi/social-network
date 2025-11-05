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

type ProfileRepository struct {
	db *pgxpool.Pool
}

const profilesTable = "profiles"

func CreateProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, user_id, first_name, second_name, birthdate, gender, biagraphy, city_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, profilesTable)
	var newProfile *models.Profile

	err := r.db.QueryRow(
		ctx,
		query,
		profile.ID,
		profile.UserId,
		profile.FirstName,
		profile.SecondName,
		profile.Birthdate,
		profile.Biography,
		profile.CityId).Scan(&newProfile)
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
	return newProfile, nil
}

func (r *ProfileRepository) GetProfileByUserId(ctx context.Context, userId uuid.UUID) (*models.Profile, error) {
	query := fmt.Sprintf(`
		SELECT *
        FROM %v
		WHERE user_id = '$1'
		LIMIT 1
	`, profilesTable)

	var profile *models.Profile

	err := r.db.QueryRow(ctx, query, userId).Scan(&profile)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get url by %s: %w", userId, err)
	}
	return profile, nil
}
