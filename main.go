package main

import (
	"context"
	"farm-backend/internal/config"
	"farm-backend/internal/db"
	"farm-backend/internal/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	router := routes.SetupRoutes(database, cfg)

	port := cfg.Port
	if port == "" {
		port = "6060"
	}
	// set up gin server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		log.Printf("Avaliable on http://localhost:%s/api", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Server stopped")
}
