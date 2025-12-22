package public

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/service"
)

type TagHandler struct {
	tagService  *service.TagService
	postService *service.PostService
}

func NewTagHandler(tagService *service.TagService, postService *service.PostService) *TagHandler {
	return &TagHandler{
		tagService:  tagService,
		postService: postService,
	}
}

// ListTags godoc
// @Summary List all tags
// @Description Get a list of all tags with post counts
// @Tags tags
// @Produce json
// @Success 200 {object} handler.Response
// @Router /api/public/tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	tags, err := h.tagService.ListTagsWithPostCount(c.Request.Context())
	if err != nil {
		handler.InternalError(c, "Failed to fetch tags")
		return
	}

	handler.Success(c, tags)
}

// GetTagPosts godoc
// @Summary Get posts by tag
// @Description Get a paginated list of posts with a tag
// @Tags tags
// @Produce json
// @Param slug path string true "Tag slug"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/public/tags/{slug}/posts [get]
func (h *TagHandler) GetTagPosts(c *gin.Context) {
	slug := c.Param("slug")

	// Get tag by slug
	tag, err := h.tagService.GetTagBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, service.ErrTagNotFound) {
			handler.NotFound(c, "Tag not found")
			return
		}
		handler.InternalError(c, "Failed to fetch tag")
		return
	}

	pagination := handler.GetPagination(c)

	posts, total, err := h.postService.ListPublishedPostsByTag(
		c.Request.Context(),
		tag.ID,
		int32(pagination.PerPage),
		int32(pagination.Offset),
	)
	if err != nil {
		handler.InternalError(c, "Failed to fetch posts")
		return
	}

	handler.SuccessWithMeta(c, posts, pagination.ToMeta(total))
}
