package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Romasmi/social-network/internal/domain/profile"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProfileRepository struct {
	writer DBQuerier
	reader DBQuerier
}

const profilesTable = "profiles"

func CreateProfileRepository(writer DBQuerier, reader DBQuerier) *ProfileRepository {
	return &ProfileRepository{writer: writer, reader: reader}
}

func (r *ProfileRepository) CreateProfile(ctx context.Context, profileModel *profile.Profile) (*profile.Profile, error) {
	const query = `
		INSERT INTO %s (id, user_id, first_name, second_name, birthdate, gender, biography, city_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, first_name, second_name, birthdate, gender, biography, city_id
	`
	sql := fmt.Sprintf(query, profilesTable)

	newProfile := &profile.Profile{}
	err := r.writer.QueryRow(
		ctx,
		sql,
		profileModel.ID,
		profileModel.UserId,
		profileModel.FirstName,
		profileModel.SecondName,
		profileModel.Birthdate,
		profileModel.Gender,
		profileModel.Biography,
		profileModel.CityId).Scan(
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

func (r *ProfileRepository) GetProfileByProfileId(ctx context.Context, profileId uuid.UUID) (*profile.Profile, error) {
	const query = `
		SELECT id, user_id, first_name, second_name, birthdate, gender, biography, city_id
        FROM %s
		WHERE id = $1
		LIMIT 1
	`
	sql := fmt.Sprintf(query, profilesTable)

	profile := &profile.Profile{}
	err := r.reader.QueryRow(ctx, sql, profileId).Scan(
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

func (r *ProfileRepository) SearchProfile(ctx context.Context, queryParams *profile.UserSearchParams) ([]*profile.Profile, error) {
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

	rows, err := r.reader.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search profiles: %w", err)
	}
	defer rows.Close()

	profiles := make([]*profile.Profile, 0)
	for rows.Next() {
		profile := &profile.Profile{}
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
