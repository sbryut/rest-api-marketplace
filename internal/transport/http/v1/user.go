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
	}
}

type userInput struct {
	Login    string `json:"login" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=8,max=64"`
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

	res, err := h.services.Users.SignIn(c.Request().Context(), service.UserInput{
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

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
