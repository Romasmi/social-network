package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CityRepository struct {
	db *pgxpool.Pool
}

const citiesTable = "cities"

func CreateCityRepository(db *pgxpool.Pool) *CityRepository {
	return &CityRepository{db: db}
}

// GetCityByName For simplicity assume that each city name is unique so it can return city by only single name
func (r *CityRepository) GetCityByName(ctx context.Context, name string) (*models.City, error) {
	const query = `
		SELECT id, name, country_code, state_province, created_at
        FROM %s
		WHERE name = $1
		LIMIT 1
	`
	sql := fmt.Sprintf(query, citiesTable)

	city := &models.City{}
	err := r.db.QueryRow(ctx, sql, name).Scan(
		&city.ID,
		&city.Name,
		&city.CountryCode,
		&city.StateProvince,
		&city.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get city by %s: %w", name, err)
	}
	return city, nil
}
