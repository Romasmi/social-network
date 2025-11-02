package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConnection struct {
	DB     *pgxpool.Pool
	Config *config.Database
}

func (c *DbConnection) Connect() error {
	pgConfig, err := pgxpool.ParseConfig(c.Config.URL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}
	pgConfig.MaxConns = int32(c.Config.MaxConnections)
	pgConfig.MinConns = int32(c.Config.MinConnections)
	pgConfig.MaxConnLifetime = utils.MinutesToNanoseconds(c.Config.MaxConnectionLifetime)
	pgConfig.MaxConnIdleTime = utils.MinutesToNanoseconds(c.Config.MaxConnectionIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c.DB, err = pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := c.Ping(); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	return nil
}

func (c *DbConnection) Close() {
	if c.DB != nil {
		c.DB.Close()
		log.Println("database connection closed")
	}
}

func (c *DbConnection) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return c.DB.Ping(ctx)
}
