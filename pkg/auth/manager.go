// Package auth provides JWT and refresh token management for authentication
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager defines methods for creating and parsing tokens
type TokenManager interface {
	NewJWTToken(userID int64, ttl time.Duration) (string, error)
	ParseJWTToken(accessToken string) (int64, error)
	NewRefreshToken() (string, error)
}

// Manager implements TokenManager using HMAC signing
type Manager struct {
	signingKey string
}

// NewManager creates a new Manager with the given signing key
func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &Manager{signingKey: signingKey}, nil
}

// TokenClaims defines custom JWT claims including user ID
type TokenClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

// NewJWTToken generates a signed JWT token with user ID and expiration time
func (m *Manager) NewJWTToken(userID int64, ttl time.Duration) (string, error) {
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.signingKey))
}

// ParseJWTToken validates a JWT token and extracts the user ID from claims
func (m *Manager) ParseJWTToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
			}
			return []byte(m.signingKey), nil
		})
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token claims: %w", err)
	}

	return claims.UserID, nil
}

// NewRefreshToken generates a secure random refresh token
func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
