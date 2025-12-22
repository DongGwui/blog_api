package main

import (
	"log"

	"github.com/ydonggwui/blog-api/internal/config"
	"github.com/ydonggwui/blog-api/internal/database"
	"github.com/ydonggwui/blog-api/internal/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to PostgreSQL
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	// Connect to Redis
	redis, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()
	log.Println("Connected to Redis")

	// Connect to MinIO
	minio, err := database.NewMinIOClient(&cfg.MinIO)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	log.Println("Connected to MinIO")

	// Setup router
	r := router.New(cfg, db, redis, minio)

	// Start server
	log.Printf("Starting server on :%s", cfg.Server.Port)
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
