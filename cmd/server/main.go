package main

import (
	"context"
	"log"
	"time"

	"github.com/ydonggwui/blog-api/internal/config"
	"github.com/ydonggwui/blog-api/internal/database"
	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/router"
	"github.com/ydonggwui/blog-api/internal/service"
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

	// Initialize sqlc queries
	queries := sqlc.New(db)

	// Connect to Redis
	redisClient, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("Connected to Redis")

	// Connect to MinIO
	minioClient, err := database.NewMinIOClient(&cfg.MinIO)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	log.Println("Connected to MinIO")

	// Seed initial admin
	if err := seedAdmin(queries, cfg); err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	// Setup router
	r := router.New(cfg, db, queries, redisClient, minioClient)

	// Start server
	log.Printf("Starting server on :%s", cfg.Server.Port)
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func seedAdmin(queries *sqlc.Queries, cfg *config.Config) error {
	if cfg.Admin.Username == "" || cfg.Admin.Password == "" {
		log.Println("Admin credentials not configured, skipping admin seed")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	authService := service.NewAuthService(queries, &cfg.JWT)
	if err := authService.EnsureAdminExists(ctx, cfg.Admin.Username, cfg.Admin.Password); err != nil {
		return err
	}

	log.Println("Admin user ensured")
	return nil
}
