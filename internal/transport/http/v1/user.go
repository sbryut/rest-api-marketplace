package v1

import (
	"errors"
	"net/http"

	"rest-api-marketplace/internal/entity"
	"rest-api-marketplace/internal/service"

	"github.com/labstack/echo/v4"
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
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	user, err := h.services.Users.SignUp(c.Request().Context(), service.UserInput{
		Login:    input.Login,
		Password: input.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserExists):
			return echo.NewHTTPError(http.StatusConflict, "user with this login already exists")
		case errors.Is(err, entity.ErrInvalidInput):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
		}
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) userSignIn(c echo.Context) error {
	var input userInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	tokens, err := h.services.Users.SignIn(c.Request().Context(), service.UserInput{
		Login:    input.Login,
		Password: input.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidCreds):
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid login or password")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to sign in")
		}
	}

	return c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *Handler) userRefresh(c echo.Context) error {
	var input refreshInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(&input); err != nil {
		return err
	}

	tokens, err := h.services.Users.RefreshTokens(c.Request().Context(), input.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound) || errors.Is(err, entity.ErrInvalidInput):
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired refresh token")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
	}

	return c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
