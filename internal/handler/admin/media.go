package admin

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/service"
)

type MediaHandler struct {
	mediaService *service.MediaService
}

func NewMediaHandler(mediaService *service.MediaService) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
	}
}

// ListMedia godoc
// @Summary List all media files
// @Description Get a paginated list of all uploaded media files
// @Tags admin/media
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} handler.Response
// @Router /api/admin/media [get]
func (h *MediaHandler) ListMedia(c *gin.Context) {
	pagination := handler.GetPagination(c)

	result, err := h.mediaService.ListMedia(
		c.Request.Context(),
		int32(pagination.PerPage),
		int32(pagination.Offset),
	)
	if err != nil {
		handler.InternalError(c, "Failed to fetch media")
		return
	}

	handler.SuccessWithMeta(c, result.Items, pagination.ToMeta(result.Total))
}

// UploadMedia godoc
// @Summary Upload a media file
// @Description Upload an image file to the server
// @Tags admin/media
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file to upload"
// @Success 201 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 413 {object} handler.ErrorResponse
// @Router /api/admin/media/upload [post]
func (h *MediaHandler) UploadMedia(c *gin.Context) {
	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		handler.BadRequest(c, "No file uploaded")
		return
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		handler.BadRequest(c, "Failed to read file")
		return
	}
	defer file.Close()

	// Get content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload file
	result, err := h.mediaService.UploadFile(
		c.Request.Context(),
		file,
		fileHeader.Filename,
		contentType,
		fileHeader.Size,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidFileType) {
			handler.BadRequest(c, "Invalid file type. Only images (JPEG, PNG, GIF, WebP, SVG) are allowed")
			return
		}
		if errors.Is(err, service.ErrFileTooLarge) {
			handler.Error(c, 413, "REQUEST_ENTITY_TOO_LARGE", "File too large. Maximum size is 10MB")
			return
		}
		handler.InternalError(c, "Failed to upload file")
		return
	}

	handler.Created(c, result)
}

// DeleteMedia godoc
// @Summary Delete a media file
// @Description Delete a media file from storage and database
// @Tags admin/media
// @Security BearerAuth
// @Param id path int true "Media ID"
// @Success 204 "No Content"
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/admin/media/{id} [delete]
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid media ID")
		return
	}

	if err := h.mediaService.DeleteMedia(c.Request.Context(), int32(id)); err != nil {
		if errors.Is(err, service.ErrMediaNotFound) {
			handler.NotFound(c, "Media not found")
			return
		}
		handler.InternalError(c, "Failed to delete media")
		return
	}

	handler.NoContent(c)
}
