package admin

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/model"
	"github.com/ydonggwui/blog-api/internal/service"
)

type TagHandler struct {
	tagService *service.TagService
}

func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// ListTags godoc
// @Summary List all tags (admin)
// @Description Get a list of all tags
// @Tags admin/tags
// @Security BearerAuth
// @Produce json
// @Success 200 {object} handler.Response
// @Router /api/admin/tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	tags, err := h.tagService.ListTags(c.Request.Context())
	if err != nil {
		handler.InternalError(c, "Failed to fetch tags")
		return
	}

	handler.Success(c, tags)
}

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag
// @Tags admin/tags
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateTagRequest true "Tag data"
// @Success 201 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /api/admin/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req model.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	tag, err := h.tagService.CreateTag(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrTagSlugExists) {
			handler.Conflict(c, "Tag slug already exists")
			return
		}
		handler.InternalError(c, "Failed to create tag")
		return
	}

	handler.Created(c, tag)
}

// UpdateTag godoc
// @Summary Update a tag
// @Description Update an existing tag
// @Tags admin/tags
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param request body model.UpdateTagRequest true "Tag data"
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /api/admin/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid tag ID")
		return
	}

	var req model.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	tag, err := h.tagService.UpdateTag(c.Request.Context(), int32(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrTagNotFound) {
			handler.NotFound(c, "Tag not found")
			return
		}
		if errors.Is(err, service.ErrTagSlugExists) {
			handler.Conflict(c, "Tag slug already exists")
			return
		}
		handler.InternalError(c, "Failed to update tag")
		return
	}

	handler.Success(c, tag)
}

// DeleteTag godoc
// @Summary Delete a tag
// @Description Delete a tag
// @Tags admin/tags
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 204 "No Content"
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/admin/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid tag ID")
		return
	}

	if err := h.tagService.DeleteTag(c.Request.Context(), int32(id)); err != nil {
		if errors.Is(err, service.ErrTagNotFound) {
			handler.NotFound(c, "Tag not found")
			return
		}
		handler.InternalError(c, "Failed to delete tag")
		return
	}

	handler.NoContent(c)
}
