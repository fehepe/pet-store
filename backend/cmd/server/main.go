package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fehepe/pet-store/backend/internal/app"
	"github.com/fehepe/pet-store/backend/internal/config"
	"github.com/fehepe/pet-store/backend/internal/server"
)

func main() {
	log.Println("Starting Pet Store Backend Server...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Environment: %s", cfg.Env)

	deps, err := app.InitializeDependencies(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer deps.Close()

	srv := server.New(deps)

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		log.Printf("GraphQL endpoint: http://localhost:%s/graphql", cfg.Port)
		if cfg.Env == "development" {
			log.Printf("Playground: http://localhost:%s/playground", cfg.Port)
		}

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown completed")
}
