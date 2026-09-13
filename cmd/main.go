package main

import (
	"log"
	"net/http"
	"os"

	"message-board/internal/repository"
	"message-board/internal/web"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is not set")
	}

	repo, err := repository.NewRedisRepository(redisURL)
	if err != nil {
		log.Fatalf("create Redis repository: %v", err)
	}

	server := web.NewRest(repo)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server.Router(),
	}

	log.Printf("HTTP server listening on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
