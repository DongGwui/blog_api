package admin

import (
	"github.com/gin-gonic/gin"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/mapper"
)

type DashboardHandler struct {
	dashboardService domainService.DashboardService
}

// NewDashboardHandlerWithCleanArch creates a new DashboardHandler with clean architecture service
func NewDashboardHandlerWithCleanArch(dashboardService domainService.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetStats godoc
// @Summary Get dashboard statistics
// @Description Get statistics for the admin dashboard
// @Tags admin/dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} handler.Response
// @Router /api/admin/dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.dashboardService.GetStats(c.Request.Context())
	if err != nil {
		handler.InternalError(c, "Failed to fetch dashboard stats")
		return
	}

	handler.Success(c, mapper.ToDashboardStatsResponse(stats))
}
