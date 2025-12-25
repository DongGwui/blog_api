package mapper

import (
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
)

// ToLoginResponse converts entity.TokenInfo to dto.LoginResponse
func ToLoginResponse(t *entity.TokenInfo) dto.LoginResponse {
	return dto.LoginResponse{
		Token:     t.Token,
		ExpiresAt: t.ExpiresAt,
	}
}

// ToAdminResponse converts entity.Admin to dto.AdminResponse
func ToAdminResponse(a *entity.Admin) dto.AdminResponse {
	return dto.AdminResponse{
		ID:        a.ID,
		Username:  a.Username,
		CreatedAt: a.CreatedAt,
	}
}
