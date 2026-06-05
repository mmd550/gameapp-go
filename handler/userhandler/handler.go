package userhandler

import (
	"encoding/json"
	"fmt"
	"gameapp/pkg/httputils"
	"gameapp/service/authservice"
	"net/http"

	"gameapp/service/userservice"
)

type Handler struct {
	userService userservice.Service
	authService *authservice.Service
}

func New(svc userservice.Service, authSvc *authservice.Service) *Handler {
	return &Handler{userService: svc, authService: authSvc}
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req userservice.RegisterRequest

	if err := decodeRequest(w, r, &req); err != nil {
		return
	}

	resp, err := handler.userService.Register(req)
	if err != nil {
		httputils.JsonError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		httputils.JsonError(w, "unexpected error", http.StatusInternalServerError)
		fmt.Printf("Register unexpected error: %v\n", err)
	}
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req userservice.LoginRequest

	if err := decodeRequest(w, r, &req); err != nil {
		return
	}

	resp, err := handler.userService.Login(req)
	if err != nil {
		httputils.JsonError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		httputils.JsonError(w, "unexpected error", http.StatusInternalServerError)
		fmt.Printf("Login unexpected error: %v\n", err)
	}
}

func (handler *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userId, authenticationError := handler.authService.AuthenticateUser(r)

	if authenticationError != nil {
		httputils.JsonError(w, authenticationError.Error(), http.StatusUnauthorized)
		return
	}

	profile, err := handler.userService.GetProfile(userservice.GetProfileRequest{
		UserId: userId,
	})

	if err != nil {
		httputils.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(profile)

	if err != nil {
		httputils.JsonError(w, "unexpected error", http.StatusInternalServerError)
		fmt.Printf("GetProfile unexpected error: %v\n", err)
	}
}

func decodeRequest(w http.ResponseWriter, r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		httputils.JsonError(w, "invalid request body", http.StatusBadRequest)
		return err
	}

	return nil
}
