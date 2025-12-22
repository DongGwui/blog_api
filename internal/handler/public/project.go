package public

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
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
// @Summary List all projects
// @Description Get a list of all projects ordered by sort_order
// @Tags projects
// @Produce json
// @Param featured query bool false "Filter by featured projects only"
// @Success 200 {object} handler.Response
// @Router /api/public/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	// Check if featured filter is requested
	if c.Query("featured") == "true" {
		projects, err := h.projectService.ListFeaturedProjects(c.Request.Context())
		if err != nil {
			handler.InternalError(c, "Failed to fetch projects")
			return
		}
		handler.Success(c, projects)
		return
	}

	projects, err := h.projectService.ListProjects(c.Request.Context())
	if err != nil {
		handler.InternalError(c, "Failed to fetch projects")
		return
	}

	handler.Success(c, projects)
}

// GetProject godoc
// @Summary Get a project by slug
// @Description Get detailed information about a project
// @Tags projects
// @Produce json
// @Param slug path string true "Project slug"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/public/projects/{slug} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	slug := c.Param("slug")

	project, err := h.projectService.GetProjectBySlug(c.Request.Context(), slug)
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
