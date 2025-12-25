package mapper

import (
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
)

// ToTagResponse converts a Tag entity to TagResponse DTO
func ToTagResponse(t *entity.Tag) *dto.TagResponse {
	resp := &dto.TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt,
	}

	if t.PostCount > 0 {
		resp.PostCount = &t.PostCount
	} else {
		zero := int64(0)
		resp.PostCount = &zero
	}

	return resp
}

// ToTagResponses converts a slice of Tag entities to TagResponse DTOs
func ToTagResponses(tags []entity.Tag) []dto.TagResponse {
	result := make([]dto.TagResponse, len(tags))
	for i := range tags {
		result[i] = *ToTagResponse(&tags[i])
	}
	return result
}

// ToCreateTagCommand converts CreateTagRequest DTO to CreateTagCommand
func ToCreateTagCommand(req *dto.CreateTagRequest) domainService.CreateTagCommand {
	return domainService.CreateTagCommand{
		Name: req.Name,
		Slug: req.Slug,
	}
}

// ToUpdateTagCommand converts UpdateTagRequest DTO to UpdateTagCommand
func ToUpdateTagCommand(req *dto.UpdateTagRequest) domainService.UpdateTagCommand {
	return domainService.UpdateTagCommand{
		Name: req.Name,
		Slug: req.Slug,
	}
}
