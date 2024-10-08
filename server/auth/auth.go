package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ValidateJWT validates the provided JWT and returns the user_id if valid
func ValidateJWT(tokenString string, secret []byte) (string, error) {
	// Parse the JWT token with the provided secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", errors.New("invalid token")
	}

	// Validate token and claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				return "", errors.New("token has expired")
			}
		} else {
			return "", errors.New("invalid expiration in token")
		}

		// Return the user_id claim
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", errors.New("user_id claim not found in token")
	}

	return "", errors.New("invalid token")
}
