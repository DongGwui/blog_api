package admin

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/model"
	"github.com/ydonggwui/blog-api/internal/service"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
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
		handler.InternalError(c, "Failed to fetch categories")
		return
	}

	handler.Success(c, categories)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags admin/categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateCategoryRequest true "Category data"
// @Success 201 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /api/admin/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrCategorySlugExists) {
			handler.Conflict(c, "Category slug already exists")
			return
		}
		handler.InternalError(c, "Failed to create category")
		return
	}

	handler.Created(c, category)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category
// @Tags admin/categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body model.UpdateCategoryRequest true "Category data"
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

	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), int32(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			handler.NotFound(c, "Category not found")
			return
		}
		if errors.Is(err, service.ErrCategorySlugExists) {
			handler.Conflict(c, "Category slug already exists")
			return
		}
		handler.InternalError(c, "Failed to update category")
		return
	}

	handler.Success(c, category)
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
		if errors.Is(err, service.ErrCategoryNotFound) {
			handler.NotFound(c, "Category not found")
			return
		}
		if errors.Is(err, service.ErrCategoryHasPosts) {
			handler.BadRequest(c, "Cannot delete category with posts")
			return
		}
		handler.InternalError(c, "Failed to delete category")
		return
	}

	handler.NoContent(c)
}
