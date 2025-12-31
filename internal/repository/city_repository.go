package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/social-network/internal/domain/city"
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
func (r *CityRepository) GetCityByName(ctx context.Context, name string) (*city.City, error) {
	const query = `
		SELECT id, name, country_code, state_province, created_at
        FROM %s
		WHERE name = $1
		LIMIT 1
	`
	sql := fmt.Sprintf(query, citiesTable)

	cityModel := &city.City{}
	err := r.db.QueryRow(ctx, sql, name).Scan(
		&cityModel.ID,
		&cityModel.Name,
		&cityModel.CountryCode,
		&cityModel.StateProvince,
		&cityModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get cityModel by %s: %w", name, err)
	}
	return cityModel, nil
}

func (r *CityRepository) CreateCity(ctx context.Context, cityModel *city.City) (*city.City, error) {
	const query = `
		INSERT INTO %s (id, name, country_code, state_province)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, country_code, state_province, created_at
	`
	sql := fmt.Sprintf(query, citiesTable)
	created := &city.City{}
	err := r.db.QueryRow(ctx, sql, cityModel.ID, cityModel.Name, cityModel.CountryCode, cityModel.StateProvince).Scan(
		&created.ID,
		&created.Name,
		&created.CountryCode,
		&created.StateProvince,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create cityModel: %w", err)
	}
	return created, nil
}
