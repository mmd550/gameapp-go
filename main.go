package main

import (
	"encoding/json"
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

	if err := postgres.RunMigrations(cfg.DB, "repository/postgres/migrations"); err != nil {
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
	mux.HandleFunc("POST /users/login", userHandler.Login)
	mux.HandleFunc("GET /health-check", healthCheckHandler)

	addr := fmt.Sprintf(":%s", cfg.HTTP.Port)
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	jsonResponse, err := json.Marshal(map[string]string{"status": "ok"})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("HealthCheckHandler: Failed to marshal response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)

	if err != nil {
		fmt.Printf("healthCheckHandler: %v", err)
	}
}
