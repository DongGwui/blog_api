package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/model"
	"github.com/ydonggwui/blog-api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary Admin login
// @Description Authenticate admin and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Router /api/admin/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			handler.Unauthorized(c, "Invalid username or password")
			return
		}
		handler.InternalError(c, "Login failed")
		return
	}

	handler.Success(c, resp)
}

// Logout godoc
// @Summary Admin logout
// @Description Logout admin (client should discard token)
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/admin/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT is stateless, so we just return success
	// Client should discard the token
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "Logged out successfully",
		},
	})
}

// Me godoc
// @Summary Get current admin info
// @Description Get the currently authenticated admin's information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.AdminResponse
// @Failure 401 {object} handler.ErrorResponse
// @Router /api/admin/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handler.Unauthorized(c, "User not authenticated")
		return
	}

	// Convert to int32 (middleware stores int64, but our DB uses int32)
	var id int32
	switch v := userID.(type) {
	case int64:
		id = int32(v)
	case int32:
		id = v
	case int:
		id = int32(v)
	default:
		handler.InternalError(c, "Invalid user ID type")
		return
	}

	admin, err := h.authService.GetAdminByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAdminNotFound) {
			handler.NotFound(c, "Admin not found")
			return
		}
		handler.InternalError(c, "Failed to get admin info")
		return
	}

	handler.Success(c, admin)
}
