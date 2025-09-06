// Package middleware provides HTTP middlewares for authentication
package middleware

import (
	"net/http"
	"strings"

	"rest-api-marketplace/pkg/auth"

	"github.com/labstack/echo/v4"
)

// AuthHeader and CtxUserID constants used for JWT authentication and context storage
const (
	AuthHeader = "Authorization"
	CtxUserID  = "user_id"
)

// JWTAuth enforces JWT authentication and sets user ID in context
func JWTAuth(tm auth.TokenManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(AuthHeader)
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization header is required"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid auth header format"})
			}
			userID, err := tm.ParseJWTToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			c.Set(CtxUserID, userID)
			return next(c)
		}
	}
}

// JWTOptionalAuth optionally parses JWT token if provided
func JWTOptionalAuth(tm auth.TokenManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(AuthHeader)
			if authHeader == "" {
				return next(c)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return next(c)
			}

			userID, err := tm.ParseJWTToken(parts[1])
			if err == nil {
				c.Set(CtxUserID, userID)
			}

			return next(c)
		}
	}
}
