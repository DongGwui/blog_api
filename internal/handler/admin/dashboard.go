package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/service"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
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

	handler.Success(c, stats)
}
