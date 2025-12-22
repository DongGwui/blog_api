package public

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/service"
)

type PostHandler struct {
	postService *service.PostService
	viewService *service.ViewService
}

func NewPostHandler(postService *service.PostService, viewService *service.ViewService) *PostHandler {
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
	if categoryID := c.Query("category"); categoryID != "" {
		var catID int32
		if _, err := parseID(categoryID, &catID); err == nil {
			posts, total, err := h.postService.ListPublishedPostsByCategory(
				c.Request.Context(),
				catID,
				int32(pagination.PerPage),
				int32(pagination.Offset),
			)
			if err != nil {
				handler.InternalError(c, "Failed to fetch posts")
				return
			}
			handler.SuccessWithMeta(c, posts, pagination.ToMeta(total))
			return
		}
	}

	// Check for tag filter
	if tagID := c.Query("tag"); tagID != "" {
		var tID int32
		if _, err := parseID(tagID, &tID); err == nil {
			posts, total, err := h.postService.ListPublishedPostsByTag(
				c.Request.Context(),
				tID,
				int32(pagination.PerPage),
				int32(pagination.Offset),
			)
			if err != nil {
				handler.InternalError(c, "Failed to fetch posts")
				return
			}
			handler.SuccessWithMeta(c, posts, pagination.ToMeta(total))
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
		handler.InternalError(c, "Failed to fetch posts")
		return
	}

	handler.SuccessWithMeta(c, posts, pagination.ToMeta(total))
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

	post, err := h.postService.GetPublishedPostBySlug(c.Request.Context(), slug)
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
		handler.InternalError(c, "Failed to search posts")
		return
	}

	handler.SuccessWithMeta(c, posts, pagination.ToMeta(total))
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
		if errors.Is(err, service.ErrPostNotFound) {
			handler.NotFound(c, "Post not found")
			return
		}
		handler.InternalError(c, "Failed to find post")
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

// Helper function to parse ID from string
func parseID(s string, id *int32) (int32, error) {
	var n int
	_, err := gin.DefaultErrorWriter.Write([]byte{})
	if err != nil {
		return 0, err
	}

	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("invalid id")
		}
		n = n*10 + int(c-'0')
	}
	*id = int32(n)
	return int32(n), nil
}
