package admin

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/model"
	"github.com/ydonggwui/blog-api/internal/service"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// ListProjects godoc
// @Summary List all projects (admin)
// @Description Get a list of all projects
// @Tags admin/projects
// @Security BearerAuth
// @Produce json
// @Success 200 {object} handler.Response
// @Router /api/admin/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	projects, err := h.projectService.ListProjects(c.Request.Context())
	if err != nil {
		handler.InternalError(c, "Failed to fetch projects")
		return
	}

	handler.Success(c, projects)
}

// GetProject godoc
// @Summary Get a project by ID (admin)
// @Description Get detailed information about a project
// @Tags admin/projects
// @Security BearerAuth
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/admin/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetProjectByID(c.Request.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			handler.NotFound(c, "Project not found")
			return
		}
		handler.InternalError(c, "Failed to fetch project")
		return
	}

	handler.Success(c, project)
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project
// @Tags admin/projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateProjectRequest true "Project data"
// @Success 201 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /api/admin/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	project, err := h.projectService.CreateProject(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrProjectSlugExists) {
			handler.Conflict(c, "Project slug already exists")
			return
		}
		handler.InternalError(c, "Failed to create project")
		return
	}

	handler.Created(c, project)
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update an existing project
// @Tags admin/projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body model.UpdateProjectRequest true "Project data"
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /api/admin/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid project ID")
		return
	}

	var req model.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), int32(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			handler.NotFound(c, "Project not found")
			return
		}
		if errors.Is(err, service.ErrProjectSlugExists) {
			handler.Conflict(c, "Project slug already exists")
			return
		}
		handler.InternalError(c, "Failed to update project")
		return
	}

	handler.Success(c, project)
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project
// @Tags admin/projects
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 204 "No Content"
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/admin/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		handler.BadRequest(c, "Invalid project ID")
		return
	}

	if err := h.projectService.DeleteProject(c.Request.Context(), int32(id)); err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			handler.NotFound(c, "Project not found")
			return
		}
		handler.InternalError(c, "Failed to delete project")
		return
	}

	handler.NoContent(c)
}

// ReorderProjects godoc
// @Summary Reorder projects
// @Description Update the sort order of multiple projects
// @Tags admin/projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.ReorderProjectsRequest true "Reorder data"
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/admin/projects/reorder [patch]
func (h *ProjectHandler) ReorderProjects(c *gin.Context) {
	var req model.ReorderProjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.projectService.ReorderProjects(c.Request.Context(), req.Orders); err != nil {
		handler.InternalError(c, "Failed to reorder projects")
		return
	}

	handler.Success(c, gin.H{"message": "Projects reordered successfully"})
}
