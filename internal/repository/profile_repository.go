package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProfileRepository struct {
	db DBQuerier
}

const profilesTable = "profiles"

func CreateProfileRepository(db DBQuerier) *ProfileRepository {
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

func (r *ProfileRepository) GetProfileByProfileId(ctx context.Context, profileId uuid.UUID) (*models.Profile, error) {
	const query = `
		SELECT id, user_id, first_name, second_name, birthdate, gender, biography, city_id
        FROM %s
		WHERE id = $1
		LIMIT 1
	`
	sql := fmt.Sprintf(query, profilesTable)

	profile := &models.Profile{}
	err := r.db.QueryRow(ctx, sql, profileId).Scan(
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

func (r *ProfileRepository) SearchProfile(ctx context.Context, queryParams *models.UserSearchParams) ([]*models.Profile, error) {
	query := `
		SELECT id, user_id, first_name, second_name, birthdate, gender, biography, city_id 
		FROM ` + profilesTable

	var conditions []string
	var args []interface{}

	if queryParams.FirstName != "" {
		conditions = append(conditions, "first_name LIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, queryParams.FirstName+"%")
	}

	if queryParams.SecondName != "" {
		conditions = append(conditions, "second_name LIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, queryParams.SecondName+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += ` 
		ORDER BY id
		LIMIT 100;
	`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search profiles: %w", err)
	}
	defer rows.Close()

	profiles := make([]*models.Profile, 0)
	for rows.Next() {
		profile := &models.Profile{}
		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan profile row: %w", err)
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating profile rows: %w", err)
	}

	return profiles, nil
}
