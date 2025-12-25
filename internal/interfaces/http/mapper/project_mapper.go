package mapper

import (
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
)

// ToCreateProjectCommand converts CreateProjectRequest to CreateProjectCommand
func ToCreateProjectCommand(req *dto.CreateProjectRequest) domainService.CreateProjectCommand {
	cmd := domainService.CreateProjectCommand{
		Title:     req.Title,
		TechStack: req.TechStack,
		Images:    req.Images,
	}

	if req.Slug != nil {
		cmd.Slug = *req.Slug
	}
	if req.Description != nil {
		cmd.Description = *req.Description
	}
	if req.Content != nil {
		cmd.Content = *req.Content
	}
	if req.DemoURL != nil {
		cmd.DemoURL = *req.DemoURL
	}
	if req.GithubURL != nil {
		cmd.GithubURL = *req.GithubURL
	}
	if req.Thumbnail != nil {
		cmd.Thumbnail = *req.Thumbnail
	}
	if req.IsFeatured != nil {
		cmd.IsFeatured = *req.IsFeatured
	}
	if req.SortOrder != nil {
		cmd.SortOrder = *req.SortOrder
	}

	return cmd
}

// ToUpdateProjectCommand converts UpdateProjectRequest to UpdateProjectCommand
func ToUpdateProjectCommand(req *dto.UpdateProjectRequest) domainService.UpdateProjectCommand {
	return domainService.UpdateProjectCommand{
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		Content:     req.Content,
		TechStack:   req.TechStack,
		DemoURL:     req.DemoURL,
		GithubURL:   req.GithubURL,
		Thumbnail:   req.Thumbnail,
		Images:      req.Images,
		IsFeatured:  req.IsFeatured,
		SortOrder:   req.SortOrder,
	}
}

// ToProjectOrderEntities converts ProjectOrderItem slice to entity.ProjectOrder slice
func ToProjectOrderEntities(orders []dto.ProjectOrderItem) []entity.ProjectOrder {
	result := make([]entity.ProjectOrder, len(orders))
	for i, o := range orders {
		result[i] = entity.ProjectOrder{
			ID:        o.ID,
			SortOrder: o.SortOrder,
		}
	}
	return result
}

// ToProjectResponse converts Project entity to ProjectResponse DTO
func ToProjectResponse(p *entity.Project) *dto.ProjectResponse {
	if p == nil {
		return nil
	}

	techStack := p.TechStack
	if techStack == nil {
		techStack = []string{}
	}

	images := p.Images
	if images == nil {
		images = []string{}
	}

	return &dto.ProjectResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Description: p.Description,
		Content:     p.Content,
		TechStack:   techStack,
		DemoURL:     p.DemoURL,
		GithubURL:   p.GithubURL,
		Thumbnail:   p.Thumbnail,
		Images:      images,
		IsFeatured:  p.IsFeatured,
		SortOrder:   p.SortOrder,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ToProjectListResponse converts Project entity to ProjectListResponse DTO
func ToProjectListResponse(p entity.Project) dto.ProjectListResponse {
	techStack := p.TechStack
	if techStack == nil {
		techStack = []string{}
	}

	return dto.ProjectListResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Description: p.Description,
		TechStack:   techStack,
		Thumbnail:   p.Thumbnail,
		IsFeatured:  p.IsFeatured,
		SortOrder:   p.SortOrder,
	}
}

// ToProjectListResponses converts a slice of Project entities to ProjectListResponse DTOs
func ToProjectListResponses(projects []entity.Project) []dto.ProjectListResponse {
	result := make([]dto.ProjectListResponse, len(projects))
	for i, p := range projects {
		result[i] = ToProjectListResponse(p)
	}
	return result
}
