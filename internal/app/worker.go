package app

import (
	"context"
	"fmt"

	"github.com/Romasmi/social-network/internal/events/consumer"
)

type Worker struct {
	app      *App
	consumer consumer.Consumer
}

func NewWorker(app *App) *Worker {
	return &Worker{app: app}
}

func (w *Worker) Run() {
	err := w.consumer.Start(context.Background())
	if err != nil {
		fmt.Errorf("start consumer error: %v", err)
		return
	}
	fmt.Print("Run worker")
}

func (w *Worker) Shutdown(ctx context.Context) error {
	return w.app.Shutdown(ctx)
}
