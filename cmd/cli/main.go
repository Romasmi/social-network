package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Romasmi/social-network/internal/cli"
)

func main() {
	basePath := os.Getenv("APP_BASE_PATH")
	if basePath == "" {
		basePath = "." // current directory
	}

	appInstance, err := cli.CreateApp(basePath)
	if err != nil {
		fmt.Printf("error while app init: %v", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := appInstance.Shutdown(ctx); err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	}

	fmt.Println("Server stopped")
}
