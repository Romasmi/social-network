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
	query := fmt.Sprintf(`
		SELECT *
        FROM %v
		WHERE name = $1
		LIMIT 1
	`, citiesTable)

	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query city by %s: %w", name, err)
	}
	defer rows.Close()

	city, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.City])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get city by %s: %w", name, err)
	}
	return &city, nil
}
