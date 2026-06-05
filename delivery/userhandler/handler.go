package userhandler

import (
	"gameapp/service/authservice"
	"net/http"

	"gameapp/service/userservice"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	userService userservice.Service
	authService *authservice.Service
}

func New(svc userservice.Service, authSvc *authservice.Service) *Handler {
	// TODO - say hi to your father 
	return &Handler{userService: svc, authService: authSvc}
}

func (handler *Handler) Register(c *echo.Context) error {
	var req userservice.RegisterRequest

	if bindErr := c.Bind(&req); bindErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bindErr.Error())
	}

	resp, err := handler.userService.Register(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (handler *Handler) Login(c *echo.Context) error {
	var req userservice.LoginRequest

	if bindErr := c.Bind(&req); bindErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bindErr.Error())
	}

	resp, err := handler.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (handler *Handler) GetProfile(c *echo.Context) error {
	userId, authenticationError := handler.authService.AuthenticateUser(c.Request())

	if authenticationError != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, authenticationError.Error())
	}

	profile, err := handler.userService.GetProfile(userservice.GetProfileRequest{
		UserId: userId,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, profile)
}
