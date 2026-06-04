package userhandler

import (
	"encoding/json"
	"fmt"
	"gameapp/pkg/httputils"
	"net/http"

	"gameapp/service/userservice"
)

type Handler struct {
	service userservice.Service
}

func New(svc userservice.Service) *Handler {
	return &Handler{service: svc}
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req userservice.RegisterRequest
	decoder := json.NewDecoder(r.Body)
	// decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		httputils.JsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := handler.service.Register(req)
	if err != nil {
		httputils.JsonError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		httputils.JsonError(w, "unexpected error", http.StatusInternalServerError)
		fmt.Printf("unexpected error: %v\n", err)
	}
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req userservice.LoginRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		httputils.JsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := handler.service.Login(req)
	if err != nil {
		httputils.JsonError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		httputils.JsonError(w, "unexpected error", http.StatusInternalServerError)
		fmt.Printf("unexpected error: %v\n", err)
	}
}
