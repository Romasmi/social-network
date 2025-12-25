package redis

import (
	"fmt"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/redis/go-redis/v9"
)

type Connection struct {
	Rdb    *redis.Client
	Config *config.Redis
}

func (c *Connection) Connect() {
	c.Rdb = redis.NewClient(&redis.Options{
		Addr:     c.Config.Host,
		Username: c.Config.Username,
		Password: c.Config.Password,
		DB:       0,
	})
}

func (c *Connection) Close() {
	if c.Rdb != nil {
		err := c.Rdb.Close()
		if err != nil {
			fmt.Printf("Erro while closing Redis connection %v\n", err)
		}
	}
}
