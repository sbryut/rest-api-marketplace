package v1

const (
	adsURL = "/ads"
	adURL  = "/ads/:id"
)

/*
type Handler struct {
}

func (h *Handler) Register(g *echo.Group) {
	g.GET("/ads", h.ListAds)
	g.POST("/ads", h.Create)
	g.GET("/ads/:id", h.Get)
	g.PUT("/ads/:id", h.Update)
	g.DELETE("/ads/:id", h.Delete)
}

func (h *Handler) ListAds(c echo.Context) error {
	return c.String(200, "ListAds")
}

func (h *Handler) Create(c echo.Context) error {
	return c.String(200, "CreateAd")
}

func (h *Handler) Get(c echo.Context) error {
	return c.String(200, "GetAd")
}

func (h *Handler) Update(c echo.Context) error {
	return c.String(200, "UpdateAd")
}

func (h *Handler) Delete(c echo.Context) error {
	return c.String(200, "DeleteAd")
}
*/
