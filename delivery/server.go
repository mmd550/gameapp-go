package httpserver

import (
	"gameapp/delivery/userhandler"
	"gameapp/service/authservice"
	"gameapp/service/userservice"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type HttpServer struct {
	engine *echo.Echo
}

type Config struct {
	UserService userservice.Service
	AuthService *authservice.Service
}

func New(config Config) HttpServer {
	e := echo.New()

	userHandler := userhandler.New(config.UserService, config.AuthService)

	e.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	e.Use(middleware.Recover())

	e.GET("/health-check", HealthCheckHandler)

	e.POST("/users/register", userHandler.Register)
	e.POST("/users/login", userHandler.Login)
	e.GET("/users/profile", userHandler.GetProfile)

	return HttpServer{e}
}

func (s HttpServer) Start() {
	if err := s.engine.Start(":8080"); err != nil {
		s.engine.Logger.Error("failed to start server", "error", err)
	}
}
