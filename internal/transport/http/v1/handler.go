package v1

import (
	"errors"
	"strconv"

	"rest-api-marketplace/internal/service"
	"rest-api-marketplace/pkg/auth"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *echo.Group) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initAdsRoutes(v1)
	}
}

func (h *Handler) parseIdFromPath(c echo.Context, param string) (int64, error) {
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
