package database

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	Master *pgxpool.Pool
	Slaves []*pgxpool.Pool
	Config *config.Database
	next   atomic.Uint64
}

func (c *Connection) Connect() error {
	var err error
	c.Master, err = c.connectToURL(c.Config.MasterURL)
	if err != nil {
		return fmt.Errorf("failed to connect to master database: %w", err)
	}

	for _, url := range c.Config.SlaveURLs {
		slave, err := c.connectToURL(url)
		if err != nil {
			return fmt.Errorf("failed to connect to slave database (%s): %w", url, err)
		}
		c.Slaves = append(c.Slaves, slave)
	}

	return nil
}

func (c *Connection) connectToURL(url string) (*pgxpool.Pool, error) {
	pgConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}
	pgConfig.MaxConns = int32(c.Config.MaxConnections)
	pgConfig.MinConns = int32(c.Config.MinConnections)
	pgConfig.MaxConnLifetime = utils.MinutesToNanoseconds(c.Config.MaxConnectionLifetime)
	pgConfig.MaxConnIdleTime = utils.MinutesToNanoseconds(c.Config.MaxConnectionIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

func (c *Connection) Close() {
	if c.Master != nil {
		c.Master.Close()
		log.Println("master database connection closed")
	}
	for i, slave := range c.Slaves {
		if slave != nil {
			slave.Close()
			log.Printf("slave database connection %d closed", i+1)
		}
	}
}

func (c *Connection) Writer() *pgxpool.Pool {
	return c.Master
}

func (c *Connection) Reader() *pgxpool.Pool {
	if len(c.Slaves) == 0 {
		return c.Master
	}
	idx := c.next.Add(1) % uint64(len(c.Slaves))
	return c.Slaves[idx]
}
