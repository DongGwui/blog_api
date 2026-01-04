package router

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ydonggwui/blog-api/internal/config"
	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	adminHandler "github.com/ydonggwui/blog-api/internal/handler/admin"
	publicHandler "github.com/ydonggwui/blog-api/internal/handler/public"
	"github.com/ydonggwui/blog-api/internal/middleware"

	// Clean Architecture imports
	appService "github.com/ydonggwui/blog-api/internal/application/service"
	postgresRepo "github.com/ydonggwui/blog-api/internal/infrastructure/persistence/postgres"
	redisRepo "github.com/ydonggwui/blog-api/internal/infrastructure/persistence/redis"
	minioStorage "github.com/ydonggwui/blog-api/internal/infrastructure/storage/minio"

	_ "github.com/ydonggwui/blog-api/docs/swagger"
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
	adminMediaHandler      *adminHandler.MediaHandler
	adminDashboardHandler  *adminHandler.DashboardHandler
}

func New(cfg *config.Config, db *sql.DB, queries *sqlc.Queries, redisClient *redis.Client, minioClient *minio.Client) *Router {
	gin.SetMode(cfg.Server.GinMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// ============================================
	// Clean Architecture Layer (All domains)
	// ============================================

	// Infrastructure Layer - Repositories
	categoryRepo := postgresRepo.NewCategoryRepository(queries)
	tagRepo := postgresRepo.NewTagRepository(queries)
	postRepo := postgresRepo.NewPostRepository(queries)
	projectRepo := postgresRepo.NewProjectRepository(queries)
	mediaRepo := postgresRepo.NewMediaRepository(queries)
	storageRepo := minioStorage.NewStorageRepository(minioClient, &cfg.MinIO)
	adminRepo := postgresRepo.NewAdminRepository(queries)
	dashboardRepo := postgresRepo.NewDashboardRepository(queries)
	viewRepo := redisRepo.NewViewRepository(redisClient)

	// Application Layer - Services (Clean Architecture)
	categoryServiceNew := appService.NewCategoryService(categoryRepo)
	tagServiceNew := appService.NewTagService(tagRepo)
	postServiceNew := appService.NewPostService(postRepo)
	projectServiceNew := appService.NewProjectService(projectRepo)
	mediaServiceNew := appService.NewMediaService(mediaRepo, storageRepo)
	authServiceNew := appService.NewAuthService(adminRepo, &cfg.JWT)
	dashboardServiceNew := appService.NewDashboardService(dashboardRepo)
	viewServiceNew := appService.NewViewService(viewRepo, postServiceNew)

	// ============================================
	// Initialize Handlers
	// ============================================

	// Auth Handler - Clean Architecture 사용
	authHandler := adminHandler.NewAuthHandlerWithCleanArch(authServiceNew)

	// Post Handlers - Clean Architecture 사용
	publicPostHandler := publicHandler.NewPostHandlerWithCleanArch(postServiceNew, viewServiceNew)
	adminPostHandler := adminHandler.NewPostHandlerWithCleanArch(postServiceNew)

	// Category Handlers - Clean Architecture 사용
	publicCategoryHandler := publicHandler.NewCategoryHandlerWithCleanArch(categoryServiceNew, postServiceNew)
	adminCategoryHandler := adminHandler.NewCategoryHandlerWithCleanArch(categoryServiceNew)

	// Tag Handlers - Clean Architecture 사용
	publicTagHandler := publicHandler.NewTagHandlerWithCleanArch(tagServiceNew, postServiceNew)
	adminTagHandler := adminHandler.NewTagHandlerWithCleanArch(tagServiceNew)

	// Project Handlers - Clean Architecture 사용
	publicProjectHandler := publicHandler.NewProjectHandlerWithCleanArch(projectServiceNew)
	adminProjectHandler := adminHandler.NewProjectHandlerWithCleanArch(projectServiceNew)

	// Media Handler - Clean Architecture 사용
	adminMediaHandler := adminHandler.NewMediaHandlerWithCleanArch(mediaServiceNew)

	// Dashboard Handler - Clean Architecture 사용
	adminDashboardHandler := adminHandler.NewDashboardHandlerWithCleanArch(dashboardServiceNew)

	r := &Router{
		engine:                engine,
		db:                    db,
		queries:               queries,
		redis:                 redisClient,
		minio:                 minioClient,
		config:                cfg,
		authHandler:           authHandler,
		publicPostHandler:     publicPostHandler,
		publicCategoryHandler: publicCategoryHandler,
		publicTagHandler:      publicTagHandler,
		publicProjectHandler:  publicProjectHandler,
		adminPostHandler:      adminPostHandler,
		adminCategoryHandler:  adminCategoryHandler,
		adminTagHandler:       adminTagHandler,
		adminProjectHandler:   adminProjectHandler,
		adminMediaHandler:     adminMediaHandler,
		adminDashboardHandler: adminDashboardHandler,
	}

	r.setupRoutes()

	return r
}

func (r *Router) setupRoutes() {
	// Swagger documentation
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
			public.GET("/categories/:slug", r.publicCategoryHandler.GetCategory)
			public.GET("/categories/:slug/posts", r.publicCategoryHandler.GetCategoryPosts)

			// Tags
			public.GET("/tags", r.publicTagHandler.ListTags)
			public.GET("/tags/:slug", r.publicTagHandler.GetTag)
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
			admin.GET("/media", r.adminMediaHandler.ListMedia)
			admin.POST("/media/upload", r.adminMediaHandler.UploadMedia)
			admin.DELETE("/media/:id", r.adminMediaHandler.DeleteMedia)

			// Dashboard
			admin.GET("/dashboard/stats", r.adminDashboardHandler.GetStats)
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
