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

// GetCityByName For simplicity assume that each city name is unique so it can return city by only single name
func (r *CityRepository) GetCityByName(ctx context.Context, name string) (*models.City, error) {
	query := fmt.Sprintf(`
		SELECT *
        FROM %v
		WHERE name = '$1'
		LIMIT 1
	`, citiesTable)

	var city *models.City

	err := r.db.QueryRow(ctx, query, name).Scan(&city)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get url by %s: %w", name, err)
	}
	return city, nil
}
