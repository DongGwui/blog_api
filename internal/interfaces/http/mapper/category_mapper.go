package mapper

import (
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
)

// ToCategoryResponse converts a Category entity to CategoryResponse DTO
func ToCategoryResponse(c *entity.Category) *dto.CategoryResponse {
	resp := &dto.CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		SortOrder: c.SortOrder,
		CreatedAt: c.CreatedAt,
	}

	if c.Description != "" {
		resp.Description = &c.Description
	}

	if c.PostCount > 0 {
		resp.PostCount = &c.PostCount
	} else {
		zero := int64(0)
		resp.PostCount = &zero
	}

	return resp
}

// ToCategoryResponses converts a slice of Category entities to CategoryResponse DTOs
func ToCategoryResponses(categories []entity.Category) []dto.CategoryResponse {
	result := make([]dto.CategoryResponse, len(categories))
	for i := range categories {
		result[i] = *ToCategoryResponse(&categories[i])
	}
	return result
}

// ToCreateCategoryCommand converts CreateCategoryRequest DTO to CreateCategoryCommand
func ToCreateCategoryCommand(req *dto.CreateCategoryRequest) domainService.CreateCategoryCommand {
	return domainService.CreateCategoryCommand{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
}

// ToUpdateCategoryCommand converts UpdateCategoryRequest DTO to UpdateCategoryCommand
func ToUpdateCategoryCommand(req *dto.UpdateCategoryRequest) domainService.UpdateCategoryCommand {
	return domainService.UpdateCategoryCommand{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
}
