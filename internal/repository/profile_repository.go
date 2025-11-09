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
	const query = `
		INSERT INTO %s (id, user_id, first_name, second_name, birthdate, gender, biography, city_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, first_name, second_name, birthdate, gender, biography, city_id
	`
	sql := fmt.Sprintf(query, profilesTable)

	newProfile := &models.Profile{}
	err := r.db.QueryRow(
		ctx,
		sql,
		profile.ID,
		profile.UserId,
		profile.FirstName,
		profile.SecondName,
		profile.Birthdate,
		profile.Gender,
		profile.Biography,
		profile.CityId).Scan(
		&newProfile.ID,
		&newProfile.UserId,
		&newProfile.FirstName,
		&newProfile.SecondName,
		&newProfile.Birthdate,
		&newProfile.Gender,
		&newProfile.Biography,
		&newProfile.CityId,
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
	return newProfile, nil
}

func (r *ProfileRepository) GetProfileByUserId(ctx context.Context, userId uuid.UUID) (*models.Profile, error) {
	const query = `
		SELECT id, user_id, first_name, second_name, birthdate, gender, biography, city_id
        FROM %s
		WHERE user_id = $1
		LIMIT 1
	`
	sql := fmt.Sprintf(query, profilesTable)

	profile := &models.Profile{}
	err := r.db.QueryRow(ctx, sql, userId).Scan(
		&profile.ID,
		&profile.UserId,
		&profile.FirstName,
		&profile.SecondName,
		&profile.Birthdate,
		&profile.Gender,
		&profile.Biography,
		&profile.CityId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get profile by user id: %w", err)
	}
	return profile, nil
}
