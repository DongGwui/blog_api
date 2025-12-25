package mapper

import (
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
)

// ToMediaResponse converts entity.Media to dto.MediaResponse
func ToMediaResponse(m *entity.Media) dto.MediaResponse {
	return dto.MediaResponse{
		ID:           m.ID,
		Filename:     m.Filename,
		OriginalName: m.OriginalName,
		Path:         m.Path,
		URL:          m.URL,
		MimeType:     m.MimeType,
		Size:         m.Size,
		Width:        m.Width,
		Height:       m.Height,
		CreatedAt:    m.CreatedAt,
	}
}

// ToMediaResponses converts a slice of entity.Media to dto.MediaResponse slice
func ToMediaResponses(media []entity.Media) []dto.MediaResponse {
	result := make([]dto.MediaResponse, len(media))
	for i, m := range media {
		result[i] = ToMediaResponse(&m)
	}
	return result
}

// ToMediaListResponse converts media entities to a list response
func ToMediaListResponse(media []entity.Media, total int64) dto.MediaListResponse {
	return dto.MediaListResponse{
		Items: ToMediaResponses(media),
		Total: total,
	}
}

// ToUploadMediaResponse converts entity.UploadedFile to dto.UploadMediaResponse
func ToUploadMediaResponse(f *entity.UploadedFile) dto.UploadMediaResponse {
	return dto.UploadMediaResponse{
		ID:           f.ID,
		Filename:     f.Filename,
		OriginalName: f.OriginalName,
		URL:          f.URL,
		MimeType:     f.MimeType,
		Size:         f.Size,
	}
}
