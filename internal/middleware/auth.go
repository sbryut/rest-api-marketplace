package middleware

import (
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"net/http"
	"rest-api-marketplace/pkg/auth"
	"strings"
)

const (
	AuthHeader = "Authorization"
	CtxUserID  = "user_id"
)

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
