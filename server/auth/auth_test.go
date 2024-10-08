package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

// TestValidateJWT tests the ValidateJWT function
func TestValidateJWT(t *testing.T) {
	// Define your secret key
	jwtSecret := []byte("test_secret")

	// Generate a valid token
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "12345",
		"exp":     time.Now().Add(time.Minute).Unix(), // Expires in 1 minute
	})
	signedValidToken, err := validToken.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to sign valid token: %v", err)
	}

	// Test 1: Valid token
	t.Run("ValidToken", func(t *testing.T) {
		userID, err := ValidateJWT(signedValidToken, jwtSecret)
		if err != nil {
			t.Errorf("Expected no error for valid token, got: %v", err)
		}
		if userID != "12345" {
			t.Errorf("Expected user_id '12345', got: %s", userID)
		}
	})

	// Generate an expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "67890",
		"exp":     time.Now().Add(-time.Minute).Unix(), // Already expired
	})
	signedExpiredToken, err := expiredToken.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to sign expired token: %v", err)
	}

	// Test 2: Expired token
	t.Run("ExpiredToken", func(t *testing.T) {
		_, err := ValidateJWT(signedExpiredToken, jwtSecret)
		if err == nil {
			t.Error("Expected error for expired token, got none")
		}
	})

	// Test 3: Invalid token (malformed)
	t.Run("InvalidToken", func(t *testing.T) {
		invalidToken := "invalid.token.string"
		_, err := ValidateJWT(invalidToken, jwtSecret)
		if err == nil {
			t.Error("Expected error for invalid token, got none")
		}
	})
}
