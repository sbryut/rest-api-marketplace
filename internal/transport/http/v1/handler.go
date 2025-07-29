package v1

import (
	"github.com/labstack/echo/v4"
	"rest-api-marketplace/internal/service"
	"rest-api-marketplace/pkg/auth"
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
	}

}
func parseIdFromPath(c *echo.Context, param string) (int64, error) {}
