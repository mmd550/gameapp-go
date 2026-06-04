package main

import (
	"fmt"
	"log"
	"net/http"

	"gameapp/config"
	"gameapp/handler/userhandler"
	"gameapp/repository/postgres"
	"gameapp/service/userservice"
)

func main() {
	cfg := config.Load()

	if err := postgres.RunMigrations(cfg.DB, "repository/migrations"); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	db, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	userSvc := userservice.New(userRepo)
	userHandler := userhandler.New(userSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users/register", userHandler.Register)
	mux.HandleFunc("GET /health-check", healthCheckHandler)

	addr := fmt.Sprintf(":%s", cfg.HTTP.Port)
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Everything is Good"))

	if err != nil {
		fmt.Printf("healthCheckHandler: %v", err)
	}
}
