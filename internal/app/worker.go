package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Romasmi/social-network/internal/events/consumer"
	"github.com/Romasmi/social-network/internal/events/events_registry"
	"github.com/Romasmi/social-network/internal/infra/kafka_client"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Worker struct {
	app      *App
	consumer consumer.Consumer
}

func NewWorker(app *App) *Worker {
	w := &Worker{app: app}
	w.Init(app.KafkaConnection)
	return w
}

func (w *Worker) Init(kafkaConn *kafka_client.Connection) {
	w.consumer = consumer.NewConsumer(
		kafkaConn,
		events_registry.NewEventRegistry(w.app.GetDB().DB, w.app.GetRedis().Rdb),
		consumer.Config{Topics: []string{"social-network-events"}},
	)
}

func (w *Worker) Run(ctx context.Context) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("Worker metrics listening on :8889")
		if err := http.ListenAndServe(":8889", nil); err != nil {
			fmt.Printf("worker metrics server error: %v\n", err)
		}
	}()

	err := w.consumer.Start(ctx)
	if err != nil {
		fmt.Errorf("start consumer error: %v", err)
		return
	}
	fmt.Println("Run worker")
}

func (w *Worker) Shutdown(ctx context.Context) error {
	return w.app.Shutdown(ctx)
}
