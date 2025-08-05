package v1

import (
	"errors"
	"net/http"
	"strconv"

	"rest-api-marketplace/internal/entity"
	"rest-api-marketplace/internal/middleware"
	"rest-api-marketplace/internal/service"

	"github.com/labstack/echo/v4"
)

func (h *Handler) initAdsRoutes(api *echo.Group) {
	ads := api.Group("/ads")
	{
		authMiddleware := middleware.JWTAuth(h.tokenManager)
		optionalAuthMiddleware := middleware.JWTOptionalAuth(h.tokenManager)
		ads.POST("", h.createAd, authMiddleware)
		ads.PUT("/:id", h.updateAd, authMiddleware)
		ads.GET("", h.listAds, optionalAuthMiddleware)
		//ads.GET("/id", h.getAdById, optionalAuthMiddleware)
		ads.DELETE("/:id", h.deleteAd, authMiddleware)

	}
}

type createAdInput struct {
	Title       string  `json:"title" validate:"required, min=1, max=100"`
	Description string  `json:"description" validate:"required, max=1000"`
	ImageURL    string  `json:"image_url" validate:"url"`
	Price       float64 `json:"price" validate:"gte=0"`
}

type updateAdInput struct {
	Title       *string  `json:"title,omitempty" validate:"required, min=1, max=100"`
	Description *string  `json:"description,omitempty" validate:"required, max=1000"`
	ImageURL    *string  `json:"image_url,omitempty" validate:"url"`
	Price       *float64 `json:"price,omitempty" validate:"gte=0"`
}

func (h *Handler) createAd(c echo.Context) error {
	userId, ok := c.Get(middleware.CtxUserID).(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
			"error": "invalid token, user not found",
		})
	}

	var input createAdInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}
	if err := c.Validate(&input); err != nil {
		return err
	}

	ad, err := h.services.Ads.Create(c.Request().Context(), service.CreateAdInput{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
	}, userId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "failed to create ad",
		})
	}

	return c.JSON(http.StatusCreated, ad)
}

func (h *Handler) updateAd(c echo.Context) error {
	userId, ok := c.Get(middleware.CtxUserID).(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
			"error": "invalid user context",
		})
	}

	adId, err := h.parseIdFromPath(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "invalid id format",
		})
	}

	var input createAdInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	updatedAd, err := h.services.Ads.Update(c.Request().Context(), adId, userId, service.UpdateAdInput{
		Title:       &input.Title,
		Description: &input.Description,
		ImageURl:    &input.ImageURL,
		Price:       &input.Price,
	})
	if err != nil {
		if errors.Is(err, entity.ErrAdNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{
				"error": "ad not found",
			})
		}
		if errors.Is(err, entity.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden, map[string]string{
				"error": "you don't have permission to update this ad",
			})
		}
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "error",
			// TODO изменить структуру ответа с ошибкой
		})
	}

	return c.JSON(http.StatusOK, updatedAd)
}

func (h *Handler) listAds(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	minPrice, _ := strconv.ParseFloat(c.QueryParam("max_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.QueryParam("max_price"), 64)

	params := entity.GetAdsQuery{
		Page:     page,
		Limit:    limit,
		SortBy:   c.QueryParam("sort_by"),
		SortDir:  c.QueryParam("sort_dir"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}

	var currentUserId *int64
	if userId, ok := c.Get(middleware.CtxUserID).(int64); ok {
		currentUserId = &userId
	}

	ads, err := h.services.Ads.GetAll(c.Request().Context(), params, currentUserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "failed to get ads",
		})
	}

	return c.JSON(http.StatusOK, ads)
}

func (h *Handler) deleteAd(c echo.Context) error {
	userId, ok := c.Get(middleware.CtxUserID).(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
			"error": "invalid user context",
		})
	}

	adId, err := h.parseIdFromPath(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "invalid id format",
		})
	}

	err = h.services.Ads.Delete(c.Request().Context(), adId, userId)
	if err != nil {
		if errors.Is(err, entity.ErrAdNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{
				"error": "ad not found",
			})
		}
		if errors.Is(err, entity.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden, map[string]string{
				"error": "you don't have permission to delete this ad",
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete this ad",
		})
	}
	return c.NoContent(http.StatusNoContent)
}
