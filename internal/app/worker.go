package app

import (
	"context"
	"fmt"
)

type Worker struct {
	app *App
}

func NewWorker(app *App) *Worker {
	return &Worker{app: app}
}

func (w *Worker) Run() {
	fmt.Print("Run worker")
}

func (w *Worker) Shutdown(ctx context.Context) error {
	return w.app.Shutdown(ctx)
}
