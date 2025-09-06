package v1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"rest-api-marketplace/internal/entity"
	"rest-api-marketplace/internal/service"
)

// initUsersRoutes registers user-related routes under /users
func (h *Handler) initUsersRoutes(api *echo.Group) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/auth/refresh", h.userRefresh)
	}
}

// userInput represents the request payload for sign-up and sign-in
type userInput struct {
	Login    string `json:"login" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// tokenResponse represents JWT access and refresh tokens returned to the client
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// refreshInput represents the request payload for refreshing tokens
type refreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// userSignUp handles user registration
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

// userSignIn handles user login and returns JWT tokens
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

// userRefresh handles refreshing JWT tokens using a valid refresh token
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
