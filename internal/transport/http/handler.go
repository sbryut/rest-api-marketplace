package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"rest-api-marketplace/internal/config"
	"rest-api-marketplace/internal/service"
	"rest-api-marketplace/internal/transport/http/v1"

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

func (h *Handler) Init(cfg *config.Config) *echo.Echo {
	e := echo.New()

	e.Use(
		middleware.Recover(),
		middleware.Logger(),
	)

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	h.initAPI(e)

	return e
}

func (h *Handler) initAPI(e *echo.Echo) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)

	api := e.Group("/api")
	handlerV1.Init(api)
}
