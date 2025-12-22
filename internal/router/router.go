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
	authHandler            *adminHandler.AuthHandler
	publicPostHandler      *publicHandler.PostHandler
	publicCategoryHandler  *publicHandler.CategoryHandler
	publicTagHandler       *publicHandler.TagHandler
	publicProjectHandler   *publicHandler.ProjectHandler
	adminPostHandler       *adminHandler.PostHandler
	adminCategoryHandler   *adminHandler.CategoryHandler
	adminTagHandler        *adminHandler.TagHandler
	adminProjectHandler    *adminHandler.ProjectHandler
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
	categoryService := service.NewCategoryService(queries)
	tagService := service.NewTagService(queries)
	projectService := service.NewProjectService(queries)

	// Initialize handlers
	authHandler := adminHandler.NewAuthHandler(authService)
	publicPostHandler := publicHandler.NewPostHandler(postService, viewService)
	publicCategoryHandler := publicHandler.NewCategoryHandler(categoryService, postService)
	publicTagHandler := publicHandler.NewTagHandler(tagService, postService)
	publicProjectHandler := publicHandler.NewProjectHandler(projectService)
	adminPostHandler := adminHandler.NewPostHandler(postService)
	adminCategoryHandler := adminHandler.NewCategoryHandler(categoryService)
	adminTagHandler := adminHandler.NewTagHandler(tagService)
	adminProjectHandler := adminHandler.NewProjectHandler(projectService)

	r := &Router{
		engine:               engine,
		db:                   db,
		queries:              queries,
		redis:                redisClient,
		minio:                minioClient,
		config:               cfg,
		authHandler:          authHandler,
		publicPostHandler:    publicPostHandler,
		publicCategoryHandler: publicCategoryHandler,
		publicTagHandler:     publicTagHandler,
		publicProjectHandler: publicProjectHandler,
		adminPostHandler:     adminPostHandler,
		adminCategoryHandler: adminCategoryHandler,
		adminTagHandler:      adminTagHandler,
		adminProjectHandler:  adminProjectHandler,
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
			public.GET("/categories", r.publicCategoryHandler.ListCategories)
			public.GET("/categories/:slug/posts", r.publicCategoryHandler.GetCategoryPosts)

			// Tags
			public.GET("/tags", r.publicTagHandler.ListTags)
			public.GET("/tags/:slug/posts", r.publicTagHandler.GetTagPosts)

			// Projects
			public.GET("/projects", r.publicProjectHandler.ListProjects)
			public.GET("/projects/:slug", r.publicProjectHandler.GetProject)
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
			admin.GET("/categories", r.adminCategoryHandler.ListCategories)
			admin.POST("/categories", r.adminCategoryHandler.CreateCategory)
			admin.PUT("/categories/:id", r.adminCategoryHandler.UpdateCategory)
			admin.DELETE("/categories/:id", r.adminCategoryHandler.DeleteCategory)

			// Tags
			admin.GET("/tags", r.adminTagHandler.ListTags)
			admin.POST("/tags", r.adminTagHandler.CreateTag)
			admin.PUT("/tags/:id", r.adminTagHandler.UpdateTag)
			admin.DELETE("/tags/:id", r.adminTagHandler.DeleteTag)

			// Projects
			admin.GET("/projects", r.adminProjectHandler.ListProjects)
			admin.GET("/projects/:id", r.adminProjectHandler.GetProject)
			admin.POST("/projects", r.adminProjectHandler.CreateProject)
			admin.PUT("/projects/:id", r.adminProjectHandler.UpdateProject)
			admin.DELETE("/projects/:id", r.adminProjectHandler.DeleteProject)
			admin.PATCH("/projects/reorder", r.adminProjectHandler.ReorderProjects)

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
