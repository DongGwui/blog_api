package admin

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/domain"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/mapper"
)

type CategoryHandler struct {
	categoryService domainService.CategoryService
}

// NewCategoryHandlerWithCleanArch creates a new CategoryHandler with clean architecture service
func NewCategoryHandlerWithCleanArch(categoryService domainService.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// ListCategories godoc
// @Summary List all categories (admin)
// @Description Get a list of all categories
// @Tags admin/categories
// @Security BearerAuth
// @Produce json
// @Success 200 {object} handler.Response
// @Router /api/admin/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		handler.InternalErrorWithLog(c, "Failed to fetch categories", err)
		return
	}
	handler.Success(c, mapper.ToCategoryResponses(categories))
}

// GetCategory godoc
// @Summary Get a category by ID (admin)
// @Description Get a single category by its ID
// @Tags admin/categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/admin/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid category ID")
		return
	}

	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), int32(id))
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			handler.NotFound(c, "Category not found")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to fetch category", err)
		return
	}

	handler.Success(c, mapper.ToCategoryResponse(category))
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags admin/categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /api/admin/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	cmd := mapper.ToCreateCategoryCommand(&req)
	category, err := h.categoryService.CreateCategory(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrCategorySlugExists) {
			handler.Conflict(c, "Category slug already exists")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to create category", err)
		return
	}

	handler.Created(c, mapper.ToCategoryResponse(category))
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category
// @Tags admin/categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category data"
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /api/admin/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid category ID")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	cmd := mapper.ToUpdateCategoryCommand(&req)
	category, err := h.categoryService.UpdateCategory(c.Request.Context(), int32(id), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			handler.NotFound(c, "Category not found")
			return
		}
		if errors.Is(err, domain.ErrCategorySlugExists) {
			handler.Conflict(c, "Category slug already exists")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to update category", err)
		return
	}

	handler.Success(c, mapper.ToCategoryResponse(category))
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category (must have no posts)
// @Tags admin/categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 204 "No Content"
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/admin/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid category ID")
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), int32(id)); err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			handler.NotFound(c, "Category not found")
			return
		}
		if errors.Is(err, domain.ErrCategoryHasPosts) {
			handler.BadRequest(c, "Cannot delete category with posts")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to delete category", err)
		return
	}

	handler.NoContent(c)
}
