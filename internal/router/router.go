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
	"github.com/ydonggwui/blog-api/internal/middleware"
)

type Router struct {
	engine *gin.Engine
	db     *sql.DB
	redis  *redis.Client
	minio  *minio.Client
	config *config.Config
}

func New(cfg *config.Config, db *sql.DB, redis *redis.Client, minio *minio.Client) *Router {
	gin.SetMode(cfg.Server.GinMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	r := &Router{
		engine: engine,
		db:     db,
		redis:  redis,
		minio:  minio,
		config: cfg,
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
			public.GET("/posts", notImplemented)
			public.GET("/posts/:slug", notImplemented)
			public.GET("/posts/search", notImplemented)
			public.POST("/posts/:slug/view", notImplemented)

			// Categories
			public.GET("/categories", notImplemented)

			// Tags
			public.GET("/tags", notImplemented)

			// Projects
			public.GET("/projects", notImplemented)
			public.GET("/projects/:slug", notImplemented)
		}

		// Admin routes (auth required)
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(r.config.JWT.Secret))
		{
			// Auth (login doesn't need auth middleware)
			// We'll handle login separately
			api.POST("/admin/auth/login", notImplemented)

			admin.GET("/auth/me", notImplemented)

			// Posts
			admin.GET("/posts", notImplemented)
			admin.GET("/posts/:id", notImplemented)
			admin.POST("/posts", notImplemented)
			admin.PUT("/posts/:id", notImplemented)
			admin.DELETE("/posts/:id", notImplemented)
			admin.PATCH("/posts/:id/publish", notImplemented)

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
