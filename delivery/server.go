package httpserver

import (
	"gameapp/config"
	"gameapp/delivery/userhandler"
	"gameapp/repository/postgres"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"gameapp/validation"
	uservalidation "gameapp/validation/user"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type HttpServer struct {
	engine *echo.Echo
}

type Config struct {
	userService   userservice.Service
	authService   *authservice.Service
	userValidator uservalidation.UserValidator
}

func New() HttpServer {
	e := echo.New()

	cfg := config.Load()

	if err := postgres.RunMigrations(cfg.DB, "repository/postgres/migrations"); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	db, err := postgres.New(cfg.DB)

	serverConfig := setupServerConfig(cfg, db)

	if err != nil {
		log.Fatalf("database: %v", err)
	}

	userHandler := userhandler.New(serverConfig.userService, serverConfig.authService, serverConfig.userValidator)

	e.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	e.Use(middleware.Recover())

	userHandler.AddUserRoutes(e, "/users")

	e.GET("/health-check", HealthCheckHandler)

	return HttpServer{e}
}

func (s HttpServer) Start() {
	if err := s.engine.Start(":8080"); err != nil {
		s.engine.Logger.Error("failed to start server", "error", err)
	}
}

func setupServerConfig(cfg config.Config, db *postgres.DB) Config {
	v := validation.New()
	userRepo := postgres.NewUserRepository(db)

	authService := authservice.New(cfg.Auth)
	userSvc := userservice.New(userRepo, authService)
	userValidator := uservalidation.New(v, userRepo)

	return Config{userService: userSvc, authService: authService, userValidator: userValidator}
}
