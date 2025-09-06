// Package http provides HTTP handlers and routes for the REST API
package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "rest-api-marketplace/docs"

	"rest-api-marketplace/internal/service"
	v1 "rest-api-marketplace/internal/transport/http/v1"
	"rest-api-marketplace/pkg/auth"
)

// @title REST API Marketplace
// @version 1.0
// @description This is a sample REST API for marketplace application.

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// Handler holds service dependencies and token manager for HTTP routes
type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

// NewHandler creates a new Handler with given services and token manager
func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

// Init sets up Echo instance, middlewares, and routes
func (h *Handler) Init() *echo.Echo {
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

// initAPI initializes API versioned routes
func (h *Handler) initAPI(e *echo.Echo) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)

	api := e.Group("/api")
	handlerV1.Init(api)
}
