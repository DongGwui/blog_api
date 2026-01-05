package public

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/domain"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/mapper"
)

type PostHandler struct {
	postService domainService.PostService
	viewService domainService.ViewService
}

// NewPostHandlerWithCleanArch creates a new PostHandler with clean architecture service
func NewPostHandlerWithCleanArch(postService domainService.PostService, viewService domainService.ViewService) *PostHandler {
	return &PostHandler{
		postService: postService,
		viewService: viewService,
	}
}

// ListPosts godoc
// @Summary List published posts
// @Description Get a paginated list of published posts
// @Tags posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param category query int false "Category ID filter"
// @Param tag query int false "Tag ID filter"
// @Success 200 {object} handler.Response
// @Router /api/public/posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	pagination := handler.GetPagination(c)

	// Check for category filter
	if categoryIDStr := c.Query("category"); categoryIDStr != "" {
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
		if err == nil {
			posts, total, err := h.postService.ListPublishedPostsByCategory(
				c.Request.Context(),
				int32(categoryID),
				int32(pagination.PerPage),
				int32(pagination.Offset),
			)
			if err != nil {
				handler.InternalErrorWithLog(c, "Failed to fetch posts", err)
				return
			}
			handler.SuccessWithMeta(c, mapper.ToPostListResponses(posts), pagination.ToMeta(total))
			return
		}
	}

	// Check for tag filter
	if tagIDStr := c.Query("tag"); tagIDStr != "" {
		tagID, err := strconv.ParseInt(tagIDStr, 10, 32)
		if err == nil {
			posts, total, err := h.postService.ListPublishedPostsByTag(
				c.Request.Context(),
				int32(tagID),
				int32(pagination.PerPage),
				int32(pagination.Offset),
			)
			if err != nil {
				handler.InternalErrorWithLog(c, "Failed to fetch posts", err)
				return
			}
			handler.SuccessWithMeta(c, mapper.ToPostListResponses(posts), pagination.ToMeta(total))
			return
		}
	}

	// Default: list all published posts
	posts, total, err := h.postService.ListPublishedPosts(
		c.Request.Context(),
		int32(pagination.PerPage),
		int32(pagination.Offset),
	)
	if err != nil {
		handler.InternalErrorWithLog(c, "Failed to fetch posts", err)
		return
	}

	handler.SuccessWithMeta(c, mapper.ToPostListResponses(posts), pagination.ToMeta(total))
}

// GetPost godoc
// @Summary Get a post by slug
// @Description Get a single published post by its slug
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/public/posts/{slug} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.postService.GetPublishedPost(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			handler.NotFound(c, "Post not found")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to fetch post", err)
		return
	}

	handler.Success(c, mapper.ToPostResponse(post))
}

// SearchPosts godoc
// @Summary Search posts
// @Description Search published posts by title and content
// @Tags posts
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/public/posts/search [get]
func (h *PostHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		handler.BadRequest(c, "Search query is required")
		return
	}

	pagination := handler.GetPagination(c)

	posts, total, err := h.postService.SearchPublishedPosts(
		c.Request.Context(),
		query,
		int32(pagination.PerPage),
		int32(pagination.Offset),
	)
	if err != nil {
		handler.InternalErrorWithLog(c, "Failed to search posts", err)
		return
	}

	handler.SuccessWithMeta(c, mapper.ToPostListResponses(posts), pagination.ToMeta(total))
}

// RecordView godoc
// @Summary Record a post view
// @Description Record a view for a post (with IP-based deduplication)
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/public/posts/{slug}/view [post]
func (h *PostHandler) RecordView(c *gin.Context) {
	slug := c.Param("slug")

	// Get post ID from slug
	postID, err := h.postService.GetPostIDBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			handler.NotFound(c, "Post not found")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to find post", err)
		return
	}

	// Record the view
	clientIP := c.ClientIP()
	isNew, err := h.viewService.RecordView(c.Request.Context(), postID, clientIP)
	if err != nil {
		// Log error but still return success
		// View counting shouldn't break the user experience
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"recorded": false,
				"message":  "View tracking temporarily unavailable",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"recorded": isNew,
		},
	})
}
