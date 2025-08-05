package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"rest-api-marketplace/internal/entity"
	"rest-api-marketplace/internal/service"
)

func (h *Handler) initUsersRoutes(api *echo.Group) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/auth/refresh", h.userRefresh)
	}
}

type userInput struct {
	Login    string `json:"login" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *Handler) userSignUp(c echo.Context) error {
	var input userInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	user, err := h.services.Users.SignUp(c.Request().Context(), service.UserInput{
		Login:    input.Login,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "user already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) userSignIn(c echo.Context) error {
	var input userInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	tokens, err := h.services.Users.SignIn(c.Request().Context(), service.UserInput{
		Login:    input.Login,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "user not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create user",
		})
	}

	return c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *Handler) userRefresh(c echo.Context) error {
	var input refreshInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input body",
		})
	}

	if err := c.Validate(&input); err != nil {
		return err
	}

	tokens, err := h.services.Users.RefreshTokens(c.Request().Context(), input.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired refresh token")
	}

	return c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
