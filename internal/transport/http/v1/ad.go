package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"rest-api-marketplace/internal/entity"
	"rest-api-marketplace/internal/middleware"
	"rest-api-marketplace/internal/service"
)

// initAdsRoutes registers all /ads endpoints with proper middlewares
func (h *Handler) initAdsRoutes(api *echo.Group) {
	ads := api.Group("/ads")
	{
		authMiddleware := middleware.JWTAuth(h.tokenManager)
		optionalAuthMiddleware := middleware.JWTOptionalAuth(h.tokenManager)
		ads.POST("", h.createAd, authMiddleware)
		ads.PUT("/:id", h.updateAd, authMiddleware)
		ads.GET("", h.listAds, optionalAuthMiddleware)
		ads.GET("/:id", h.getAdByID, optionalAuthMiddleware)
		ads.DELETE("/:id", h.deleteAd, authMiddleware)
	}
}

// createAdInput defines input structure for creating a new ad
type createAdInput struct {
	Title       string  `json:"title" validate:"required,min=1,max=100"`
	Description string  `json:"description" validate:"required,max=1000"`
	ImageURL    string  `json:"image_url" validate:"url"`
	Price       float64 `json:"price" validate:"gte=0"`
}

// updateAdInput defines input structure for updating an ad
type updateAdInput struct {
	Title       *string  `json:"title,omitempty" validate:"required,min=1,max=100"`
	Description *string  `json:"description,omitempty" validate:"required,max=1000"`
	ImageURL    *string  `json:"image_url,omitempty" validate:"url"`
	Price       *float64 `json:"price,omitempty" validate:"gte=0"`
}

// @Summary Create Ad
// @Description Create a new advertisement
// @Tags ads
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param ad body createAdInput true "Ad creation details"
// @Success 201 {object} entity.Ad
// @Failure 400 {object} error "Invalid request body or input"
// @Failure 409 {object} error "Unauthorized"
// @Failure 500 {object} error "Failed to create ad"
// @Router /api/v1/ads [post]
// createAd handles POST /ads to create a new advertisement
func (h *Handler) createAd(c echo.Context) error {
	userID, ok := c.Get(middleware.CtxUserID).(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user context")
	}

	var input createAdInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(&input); err != nil {
		return err
	}

	ad, err := h.services.Ads.Create(c.Request().Context(), service.CreateAdInput{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
	}, userID)

	if err != nil {
		if errors.Is(err, entity.ErrInvalidInput) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create ad")
	}

	return c.JSON(http.StatusCreated, ad)
}

// @Summary Update Ad
// @Description Update an existing advertisement
// @Tags ads
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path int64 true "Ad ID"
// @Param ad body updateAdInput true "Ad update details"
// @Success 200 {object} entity.Ad
// @Failure 400 {object} error "Invalid request body or input"
// @Failure 401 {object} error "Unauthorized"
// @Failure 403 {object} error "Forbidden"
// @Failure 404 {object} error "Ad not found"
// @Failure 500 {object} error "Failed to update ad"
// @Router /api/v1/ads/{id} [put]
// updateAd handles PUT /ads/:id to update an existing advertisement
func (h *Handler) updateAd(c echo.Context) error {
	userID, ok := c.Get(middleware.CtxUserID).(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user context")
	}

	adID, err := h.parseIDFromPath(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var input updateAdInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	updatedAd, err := h.services.Ads.Update(c.Request().Context(), adID, userID, service.UpdateAdInput{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
	})
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAdNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "ad not found")
		case errors.Is(err, entity.ErrForbidden):
			return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to update this ad")
		case errors.Is(err, entity.ErrInvalidInput):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update the ad")
		}
	}
	return c.JSON(http.StatusOK, updatedAd)
}

// @Summary List Ads
// @Description Retrieve a paginated list of advertisements
// @Tags ads
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param min_price query float64 false "Minimum price"
// @Param max_price query float64 false "Maximum price"
// @Param sort_by query string false "Sort by field (e.g., price, created_at)"
// @Param sort_dir query string false "Sort direction (asc or desc)"
// @Success 200 {array} entity.Ad
// @Failure 400 {object} error "Incorrect price parameter"
// @Failure 500 {object} error "Failed to get ads"
// @Router /api/v1/ads [get]
// listAds handles GET /ads to retrieve a paginated list of advertisements
func (h *Handler) listAds(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	var minPrice, maxPrice float64
	if mp := c.QueryParam("min_price"); mp != "" {
		val, err := strconv.ParseFloat(mp, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "incorrect min price")
		}
		minPrice = val
	}
	if mp := c.QueryParam("max_price"); mp != "" {
		val, err := strconv.ParseFloat(mp, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "incorrect max price")
		}
		maxPrice = val
	}

	params := entity.GetAdsQuery{
		Page:     page,
		Limit:    limit,
		SortBy:   c.QueryParam("sort_by"),
		SortDir:  c.QueryParam("sort_dir"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}

	var currentUserID *int64
	if val := c.Get(middleware.CtxUserID); val != nil {
		switch v := val.(type) {
		case int64:
			currentUserID = &v
		case float64:
			tmp := int64(v)
			currentUserID = &tmp
		}
	}

	ads, err := h.services.Ads.GetAll(c.Request().Context(), params, currentUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get ads")
	}

	return c.JSON(http.StatusOK, ads)
}

// @Summary Get Ad by ID
// @Description Retrieve a single advertisement by its ID
// @Tags ads
// @Produce json
// @Param id path int64 true "Ad ID"
// @Success 200 {object} entity.Ad
// @Failure 400 {object} error "Invalid ad ID"
// @Failure 404 {object} error "Ad not found"
// @Failure 500 {object} error "Failed to get ad"
// @Router /api/v1/ads/{id} [get]
// getAdByID handles GET /ads/:id to retrieve a single advertisement by ID
func (h *Handler) getAdByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ad id")
	}

	var currentUserID *int64
	if userID, ok := c.Get(middleware.CtxUserID).(int64); ok {
		currentUserID = &userID
	}

	ad, err := h.services.Ads.GetByIDWithAuthor(c.Request().Context(), id, currentUserID)
	if err != nil {
		if errors.Is(err, entity.ErrAdNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "ad not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get ad")
	}

	return c.JSON(http.StatusOK, ad)
}

// @Summary Delete Ad
// @Description Delete an advertisement by its ID
// @Tags ads
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path int64 true "Ad ID"
// @Success 204 "No content"
// @Failure 400 {object} error "Invalid ad ID"
// @Failure 401 {object} error "Unauthorized"
// @Failure 403 {object} error "Forbidden"
// @Failure 404 {object} error "Ad not found"
// @Failure 500 {object} error "Failed to delete ad"
// @Router /api/v1/ads/{id} [delete]
// deleteAd handles DELETE /ads/:id to remove an advertisement by ID
func (h *Handler) deleteAd(c echo.Context) error {
	userID, ok := c.Get(middleware.CtxUserID).(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user context")
	}

	adID, err := h.parseIDFromPath(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.services.Ads.Delete(c.Request().Context(), adID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAdNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "ad not found")
		case errors.Is(err, entity.ErrForbidden):
			return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to delete this ad")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete this ad")
		}
	}
	return c.NoContent(http.StatusNoContent)
}
