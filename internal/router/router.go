package router

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/ydonggwui/blog-api/internal/config"
	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	adminHandler "github.com/ydonggwui/blog-api/internal/handler/admin"
	publicHandler "github.com/ydonggwui/blog-api/internal/handler/public"
	"github.com/ydonggwui/blog-api/internal/middleware"
	"github.com/ydonggwui/blog-api/internal/service"
)

type Router struct {
	engine  *gin.Engine
	db      *sql.DB
	queries *sqlc.Queries
	redis   *redis.Client
	minio   *minio.Client
	config  *config.Config

	// Handlers
	authHandler       *adminHandler.AuthHandler
	publicPostHandler *publicHandler.PostHandler
	adminPostHandler  *adminHandler.PostHandler
}

func New(cfg *config.Config, db *sql.DB, queries *sqlc.Queries, redisClient *redis.Client, minioClient *minio.Client) *Router {
	gin.SetMode(cfg.Server.GinMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// Initialize services
	authService := service.NewAuthService(queries, &cfg.JWT)
	postService := service.NewPostService(queries, db)
	viewService := service.NewViewService(redisClient, postService)

	// Initialize handlers
	authHandler := adminHandler.NewAuthHandler(authService)
	publicPostHandler := publicHandler.NewPostHandler(postService, viewService)
	adminPostHandler := adminHandler.NewPostHandler(postService)

	r := &Router{
		engine:            engine,
		db:                db,
		queries:           queries,
		redis:             redisClient,
		minio:             minioClient,
		config:            cfg,
		authHandler:       authHandler,
		publicPostHandler: publicPostHandler,
		adminPostHandler:  adminPostHandler,
	}

	r.setupRoutes()

	return r
}

func (r *Router) setupRoutes() {
	api := r.engine.Group("/api")
	{
		// Health check
		api.GET("/health", r.healthCheck)

		// Public routes (no auth required)
		public := api.Group("/public")
		{
			// Posts
			public.GET("/posts", r.publicPostHandler.ListPosts)
			public.GET("/posts/search", r.publicPostHandler.SearchPosts)
			public.GET("/posts/:slug", r.publicPostHandler.GetPost)
			public.POST("/posts/:slug/view", r.publicPostHandler.RecordView)

			// Categories
			public.GET("/categories", notImplemented)

			// Tags
			public.GET("/tags", notImplemented)

			// Projects
			public.GET("/projects", notImplemented)
			public.GET("/projects/:slug", notImplemented)
		}

		// Admin auth routes (no auth required for login)
		adminAuth := api.Group("/admin/auth")
		{
			adminAuth.POST("/login", r.authHandler.Login)
			adminAuth.POST("/logout", r.authHandler.Logout)
		}

		// Admin routes (auth required)
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(r.config.JWT.Secret))
		{
			// Auth
			admin.GET("/auth/me", r.authHandler.Me)

			// Posts
			admin.GET("/posts", r.adminPostHandler.ListPosts)
			admin.GET("/posts/:id", r.adminPostHandler.GetPost)
			admin.POST("/posts", r.adminPostHandler.CreatePost)
			admin.PUT("/posts/:id", r.adminPostHandler.UpdatePost)
			admin.DELETE("/posts/:id", r.adminPostHandler.DeletePost)
			admin.PATCH("/posts/:id/publish", r.adminPostHandler.PublishPost)

			// Categories
			admin.GET("/categories", notImplemented)
			admin.POST("/categories", notImplemented)
			admin.PUT("/categories/:id", notImplemented)
			admin.DELETE("/categories/:id", notImplemented)

			// Tags
			admin.GET("/tags", notImplemented)
			admin.POST("/tags", notImplemented)
			admin.PUT("/tags/:id", notImplemented)
			admin.DELETE("/tags/:id", notImplemented)

			// Projects
			admin.GET("/projects", notImplemented)
			admin.POST("/projects", notImplemented)
			admin.PUT("/projects/:id", notImplemented)
			admin.DELETE("/projects/:id", notImplemented)
			admin.PATCH("/projects/reorder", notImplemented)

			// Media
			admin.GET("/media", notImplemented)
			admin.POST("/media/upload", notImplemented)
			admin.DELETE("/media/:id", notImplemented)

			// Dashboard
			admin.GET("/dashboard/stats", notImplemented)
		}
	}
}

func (r *Router) healthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	status := "ok"
	checks := make(map[string]string)

	// Check PostgreSQL
	if err := r.db.PingContext(ctx); err != nil {
		status = "degraded"
		checks["postgres"] = "error: " + err.Error()
	} else {
		checks["postgres"] = "ok"
	}

	// Check Redis
	if err := r.redis.Ping(ctx).Err(); err != nil {
		status = "degraded"
		checks["redis"] = "error: " + err.Error()
	} else {
		checks["redis"] = "ok"
	}

	// Check MinIO
	if _, err := r.minio.BucketExists(ctx, r.config.MinIO.Bucket); err != nil {
		status = "degraded"
		checks["minio"] = "error: " + err.Error()
	} else {
		checks["minio"] = "ok"
	}

	httpStatus := http.StatusOK
	if status != "ok" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status": status,
		"checks": checks,
	})
}

func (r *Router) Run() error {
	return r.engine.Run(":" + r.config.Server.Port)
}

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "This endpoint is not implemented yet",
		},
	})
}
