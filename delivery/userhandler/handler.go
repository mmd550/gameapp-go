package userhandler

import (
	"gameapp/dto"
	"gameapp/pkg/httpmessage"
	"gameapp/service/authservice"
	"gameapp/validation"
	uservalidation "gameapp/validation/user"
	"net/http"

	"gameapp/service/userservice"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	userService userservice.Service
	authService *authservice.Service
	validator   uservalidation.UserValidator
}

func New(svc userservice.Service, authSvc *authservice.Service, validator uservalidation.UserValidator) *Handler {
	return &Handler{userService: svc, authService: authSvc, validator: validator}
}

func (handler *Handler) Register(c *echo.Context) error {
	var req dto.RegisterRequest

	if bindErr := c.Bind(&req); bindErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bindErr.Error())
	}

	validationError, filedErrors := handler.validator.ValidateRegisterRequest(req)

	if validationError != nil {
		message, code := httpmessage.Error(validationError)

		return c.JSON(code, validation.ValidationError{
			Message: message,
			FieldErrors: filedErrors,
		})
	}

	resp, err := handler.userService.Register(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (handler *Handler) Login(c *echo.Context) error {
	var req dto.LoginRequest

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

	profile, err := handler.userService.GetProfile(dto.GetProfileRequest{
		UserId: userId,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, profile)
}
