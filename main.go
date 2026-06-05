package main

import (
	httpserver "gameapp/delivery"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"log"

	"gameapp/config"
	"gameapp/repository/postgres"
)

func main() {
	cfg := config.Load()

	if err := postgres.RunMigrations(cfg.DB, "repository/postgres/migrations"); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	db, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	serverConfig := setupServerConfig(cfg, db)

	server := httpserver.New(serverConfig)

	server.Start()
}

func setupServerConfig(cfg config.Config, db *postgres.DB) httpserver.Config {
	userRepo := postgres.NewUserRepository(db)

	authService := authservice.New(cfg.Auth)
	userSvc := userservice.New(userRepo, authService)

	return httpserver.Config{UserService: userSvc, AuthService: authService}
}
