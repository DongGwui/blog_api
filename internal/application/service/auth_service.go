package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/ydonggwui/blog-api/internal/config"
	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
)

// jwtClaims represents internal JWT claims structure
type jwtClaims struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type authService struct {
	adminRepo repository.AdminRepository
	jwtConfig *config.JWTConfig
}

func NewAuthService(adminRepo repository.AdminRepository, jwtConfig *config.JWTConfig) domainService.AuthService {
	return &authService{
		adminRepo: adminRepo,
		jwtConfig: jwtConfig,
	}
}

func (s *authService) Login(ctx context.Context, cmd domainService.LoginCommand) (*entity.TokenInfo, error) {
	admin, err := s.adminRepo.FindByUsername(ctx, cmd.Username)
	if err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := comparePassword(admin.Password, cmd.Password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, expiresAt, err := s.generateToken(admin.ID, admin.Username)
	if err != nil {
		return nil, err
	}

	return &entity.TokenInfo{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) GetAdminByID(ctx context.Context, id int32) (*entity.Admin, error) {
	return s.adminRepo.FindByID(ctx, id)
}

func (s *authService) ValidateToken(tokenString string) (*entity.Claims, error) {
	claims := &jwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &entity.Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
	}, nil
}

func (s *authService) EnsureAdminExists(ctx context.Context, username, password string) error {
	_, err := s.adminRepo.FindByUsername(ctx, username)
	if err == nil {
		// Admin already exists
		return nil
	}

	if !errors.Is(err, domain.ErrAdminNotFound) {
		return err
	}

	// Create initial admin
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.adminRepo.Create(ctx, &entity.Admin{
		Username: username,
		Password: hashedPassword,
	})

	return err
}

// generateToken creates a new JWT token
func (s *authService) generateToken(userID int32, username string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.jwtConfig.Expiry)

	claims := &jwtClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// comparePassword compares a hashed password with a plain password
func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
