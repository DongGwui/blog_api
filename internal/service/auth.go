package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/ydonggwui/blog-api/internal/config"
	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/model"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAdminNotFound      = errors.New("admin not found")
)

type AuthService struct {
	queries *sqlc.Queries
	config  *config.JWTConfig
}

func NewAuthService(queries *sqlc.Queries, cfg *config.JWTConfig) *AuthService {
	return &AuthService{
		queries: queries,
		config:  cfg,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login authenticates an admin and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	admin, err := s.queries.GetAdminByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := ComparePassword(admin.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, expiresAt, err := s.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// GenerateToken creates a new JWT token
func (s *AuthService) GenerateToken(userID int32, username string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.config.Expiry)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetAdminByID returns an admin by ID
func (s *AuthService) GetAdminByID(ctx context.Context, id int32) (*model.AdminResponse, error) {
	admin, err := s.queries.GetAdminByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAdminNotFound
		}
		return nil, err
	}

	return &model.AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		CreatedAt: admin.CreatedAt.Time,
	}, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword compares a hashed password with a plain password
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// EnsureAdminExists creates the initial admin if it doesn't exist
func (s *AuthService) EnsureAdminExists(ctx context.Context, username, password string) error {
	_, err := s.queries.GetAdminByUsername(ctx, username)
	if err == nil {
		// Admin already exists
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// Create initial admin
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.queries.CreateAdmin(ctx, sqlc.CreateAdminParams{
		Username: username,
		Password: hashedPassword,
	})

	return err
}
