package service

import (
	"testing"
	"time"

	"github.com/ydonggwui/blog-api/internal/config"
)

func TestAuthService_GenerateToken(t *testing.T) {
	jwtConfig := &config.JWTConfig{
		Secret: "test-secret-key-at-least-32-chars",
		Expiry: 24 * time.Hour,
	}

	authService := &AuthService{
		config: jwtConfig,
	}

	// Test token generation
	token, expiresAt, err := authService.GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Token should not be empty")
	}

	if expiresAt.IsZero() {
		t.Fatal("ExpiresAt should not be zero")
	}

	t.Logf("Generated token: %s...", token[:50])
}

func TestAuthService_ValidateToken(t *testing.T) {
	jwtConfig := &config.JWTConfig{
		Secret: "test-secret-key-at-least-32-chars",
		Expiry: 24 * time.Hour,
	}

	authService := &AuthService{
		config: jwtConfig,
	}

	// Generate a token first
	token, _, err := authService.GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	claims, err := authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", claims.UserID)
	}

	if claims.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", claims.Username)
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	jwtConfig := &config.JWTConfig{
		Secret: "test-secret-key-at-least-32-chars",
		Expiry: 24 * time.Hour,
	}

	authService := &AuthService{
		config: jwtConfig,
	}

	// Test with invalid token
	_, err := authService.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPassword == "" {
		t.Fatal("Hashed password should not be empty")
	}

	if hashedPassword == password {
		t.Fatal("Hashed password should not equal plain password")
	}
}

func TestComparePassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test with correct password
	if err := ComparePassword(hashedPassword, password); err != nil {
		t.Fatal("ComparePassword should not return error for correct password")
	}

	// Test with incorrect password
	if err := ComparePassword(hashedPassword, "wrongpassword"); err == nil {
		t.Fatal("ComparePassword should return error for incorrect password")
	}
}
