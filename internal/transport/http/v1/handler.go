// Package v1 provides version 1 HTTP handlers for REST API
package v1

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"

	"rest-api-marketplace/internal/service"
	"rest-api-marketplace/pkg/auth"
)

// Handler holds services and token manager to handle HTTP requests
type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

// NewHandler creates a new HTTP handler with given services and token manager
func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

// Init registers version 1 API routes under /v1
func (h *Handler) Init(api *echo.Group) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initAdsRoutes(v1)
	}
}

// parseIDFromPath extracts and validates an integer ID parameter from the URL path
func (h *Handler) parseIDFromPath(c echo.Context, param string) (int64, error) {
	idStr := c.Param(param)
	if idStr == "" {
		return 0, errors.New("empty id param")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id param")
	}

	return id, nil
}
