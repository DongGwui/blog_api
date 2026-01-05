package admin

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/domain"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/mapper"
)

type MediaHandler struct {
	mediaService domainService.MediaService
}

// NewMediaHandlerWithCleanArch creates a new MediaHandler with clean architecture service
func NewMediaHandlerWithCleanArch(mediaService domainService.MediaService) *MediaHandler {
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

	media, total, err := h.mediaService.ListMedia(
		c.Request.Context(),
		int32(pagination.PerPage),
		int32(pagination.Offset),
	)
	if err != nil {
		handler.InternalErrorWithLog(c, "Failed to fetch media", err)
		return
	}

	handler.SuccessWithMeta(c, mapper.ToMediaResponses(media), pagination.ToMeta(total))
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
	result, err := h.mediaService.UploadMedia(
		c.Request.Context(),
		domainService.UploadMediaCommand{
			File:         file,
			OriginalName: fileHeader.Filename,
			MimeType:     contentType,
			Size:         fileHeader.Size,
		},
	)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidFileType) {
			handler.BadRequest(c, "Invalid file type. Only images (JPEG, PNG, GIF, WebP, SVG) are allowed")
			return
		}
		if errors.Is(err, domain.ErrFileTooLarge) {
			handler.Error(c, 413, "REQUEST_ENTITY_TOO_LARGE", "File too large. Maximum size is 10MB")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to upload file", err)
		return
	}

	handler.Created(c, mapper.ToUploadMediaResponse(result))
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
		if errors.Is(err, domain.ErrMediaNotFound) {
			handler.NotFound(c, "Media not found")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to delete media", err)
		return
	}

	handler.NoContent(c)
}
