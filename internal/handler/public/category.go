package public

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/domain"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/mapper"
)

type CategoryHandler struct {
	categoryService domainService.CategoryService
	postService     domainService.PostService
}

// NewCategoryHandlerWithCleanArch creates a new CategoryHandler with clean architecture
func NewCategoryHandlerWithCleanArch(categoryService domainService.CategoryService, postService domainService.PostService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		postService:     postService,
	}
}

// ListCategories godoc
// @Summary List all categories
// @Description Get a list of all categories with post counts
// @Tags categories
// @Produce json
// @Success 200 {object} handler.Response
// @Router /api/public/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		handler.InternalErrorWithLog(c, "Failed to fetch categories", err)
		return
	}

	handler.Success(c, mapper.ToCategoryResponses(categories))
}

// GetCategoryPosts godoc
// @Summary Get posts by category
// @Description Get a paginated list of posts in a category
// @Tags categories
// @Produce json
// @Param slug path string true "Category slug"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/public/categories/{slug}/posts [get]
func (h *CategoryHandler) GetCategoryPosts(c *gin.Context) {
	slug := c.Param("slug")

	// Get category by slug (Clean Architecture)
	category, err := h.categoryService.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			handler.NotFound(c, "Category not found")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to fetch category", err)
		return
	}

	pagination := handler.GetPagination(c)

	// Get posts (Clean Architecture)
	posts, total, err := h.postService.ListPublishedPostsByCategory(
		c.Request.Context(),
		category.ID,
		int32(pagination.PerPage),
		int32(pagination.Offset),
	)
	if err != nil {
		handler.InternalErrorWithLog(c, "Failed to fetch posts", err)
		return
	}

	handler.SuccessWithMeta(c, mapper.ToPostListResponses(posts), pagination.ToMeta(total))
}
