package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ydonggwui/blog-api/internal/domain"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/handler"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/mapper"
)

type AuthHandler struct {
	authService domainService.AuthService
}

// NewAuthHandlerWithCleanArch creates a new AuthHandler with clean architecture service
func NewAuthHandlerWithCleanArch(authService domainService.AuthService) *AuthHandler {
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
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Router /api/admin/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.BadRequest(c, "Invalid request body")
		return
	}

	tokenInfo, err := h.authService.Login(c.Request.Context(), domainService.LoginCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			handler.Unauthorized(c, "Invalid username or password")
			return
		}
		handler.InternalErrorWithLog(c, "Login failed", err)
		return
	}

	handler.Success(c, mapper.ToLoginResponse(tokenInfo))
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
// @Success 200 {object} dto.AdminResponse
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
		if errors.Is(err, domain.ErrAdminNotFound) {
			handler.NotFound(c, "Admin not found")
			return
		}
		handler.InternalErrorWithLog(c, "Failed to get admin info", err)
		return
	}

	handler.Success(c, mapper.ToAdminResponse(admin))
}
