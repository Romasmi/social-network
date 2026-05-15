package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Romasmi/social-network/internal/config"
	"github.com/Romasmi/social-network/internal/events/publisher"
	"github.com/Romasmi/social-network/internal/infra/database"
	"github.com/Romasmi/social-network/internal/infra/kafka_client"
	"github.com/Romasmi/social-network/internal/infra/redis_client"
)

type App struct {
	DbConn          *database.Connection
	RedisConn       *redis_client.Connection
	KafkaConnection *kafka_client.Connection
	Config          *config.Config
	Publisher       publisher.Publisher
	server          *http.Server
}

func (a *App) GetWriter() *database.Connection {
	return a.DbConn
}

func (a *App) GetRedis() *redis_client.Connection {
	return a.RedisConn
}

func (a *App) GetPublisher() publisher.Publisher {
	return a.Publisher
}

func NewApp(configPath string) (*App, error) {
	app := &App{}
	return app, app.init(configPath)
}

func (a *App) init(configPath string) error {
	envConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading Config: %v\n", err)
	}
	a.Config = envConfig

	dbConn := &database.Connection{Config: &envConfig.Database}
	if err = dbConn.Connect(); err != nil {
		return fmt.Errorf("error connecting to DB: %v\n", err)
	}
	a.DbConn = dbConn

	redisConn := &redis_client.Connection{Config: &envConfig.Redis}
	redisConn.Connect()
	a.RedisConn = redisConn

	kafkaConn, err := kafka_client.CreateKafkaConnection(&a.Config.Kafka)
	if err != nil {
		return fmt.Errorf("error connecting to Kafka: %v\n", err)
	}
	a.KafkaConnection = kafkaConn

	a.Publisher = publisher.NewPublisher(a.KafkaConnection, publisher.PublisherConfig{})

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var shutdownErr error

	if a.server != nil {
		fmt.Println("Shutting down HTTP server...")
		if err := a.server.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("server shutdown error: %w", err)
			fmt.Printf("HTTP server shutdown error: %v\n", err)
		}
	}

	if a.DbConn != nil {
		fmt.Println("Closing database connections...")
		a.DbConn.Close()
	}

	if a.Publisher.Flush(time.Minute*3) != nil {
		shutdownErr = fmt.Errorf("error flushing events: %w", a.Publisher.Flush(time.Minute*3))
	}

	a.KafkaConnection.Close()

	fmt.Println("Cleanup completed")
	return shutdownErr
}
