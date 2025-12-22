package admin

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/model"
	"github.com/ydonggwui/blog-api/internal/service"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// ListPosts godoc
// @Summary List all posts (admin)
// @Description Get a paginated list of all posts with optional status filter
// @Tags admin/posts
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param status query string false "Filter by status (draft, published)"
// @Success 200 {object} handler.Response
// @Router /api/admin/posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	pagination := handler.GetPagination(c)
	status := c.Query("status")

	var posts []model.PostListResponse
	var total int64
	var err error

	if status != "" {
		posts, total, err = h.postService.ListPostsByStatus(
			c.Request.Context(),
			status,
			int32(pagination.PerPage),
			int32(pagination.Offset),
		)
	} else {
		posts, total, err = h.postService.ListAllPosts(
			c.Request.Context(),
			int32(pagination.PerPage),
			int32(pagination.Offset),
		)
	}

	if err != nil {
		handler.InternalError(c, "Failed to fetch posts")
		return
	}

	handler.SuccessWithMeta(c, posts, pagination.ToMeta(total))
}

// GetPost godoc
// @Summary Get a post by ID (admin)
// @Description Get a single post by its ID
// @Tags admin/posts
// @Security BearerAuth
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/admin/posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid post ID")
		return
	}

	post, err := h.postService.GetPostByID(c.Request.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			handler.NotFound(c, "Post not found")
			return
		}
		handler.InternalError(c, "Failed to fetch post")
		return
	}

	handler.Success(c, post)
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new blog post
// @Tags admin/posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreatePostRequest true "Post data"
// @Success 201 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /api/admin/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	post, err := h.postService.CreatePost(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrSlugExists) {
			handler.Conflict(c, "Slug already exists")
			return
		}
		handler.InternalError(c, "Failed to create post")
		return
	}

	handler.Created(c, post)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update an existing blog post
// @Tags admin/posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param request body model.UpdatePostRequest true "Post data"
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /api/admin/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid post ID")
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	post, err := h.postService.UpdatePost(c.Request.Context(), int32(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			handler.NotFound(c, "Post not found")
			return
		}
		if errors.Is(err, service.ErrSlugExists) {
			handler.Conflict(c, "Slug already exists")
			return
		}
		handler.InternalError(c, "Failed to update post")
		return
	}

	handler.Success(c, post)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a blog post
// @Tags admin/posts
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 204 "No Content"
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/admin/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid post ID")
		return
	}

	if err := h.postService.DeletePost(c.Request.Context(), int32(id)); err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			handler.NotFound(c, "Post not found")
			return
		}
		handler.InternalError(c, "Failed to delete post")
		return
	}

	handler.NoContent(c)
}

// PublishPost godoc
// @Summary Publish or unpublish a post
// @Description Change the publish status of a post
// @Tags admin/posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param request body model.PublishRequest true "Publish status"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/admin/posts/{id}/publish [patch]
func (h *PostHandler) PublishPost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid post ID")
		return
	}

	var req model.PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	post, err := h.postService.PublishPost(c.Request.Context(), int32(id), req.Publish)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			handler.NotFound(c, "Post not found")
			return
		}
		handler.InternalError(c, "Failed to update publish status")
		return
	}

	handler.Success(c, post)
}
