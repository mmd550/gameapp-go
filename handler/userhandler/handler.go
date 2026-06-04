package userhandler

import (
	"encoding/json"
	"gameapp/pkg/httputils"
	"net/http"

	"gameapp/service/userservice"
)

type Handler struct {
	svc userservice.Service
}

func New(svc userservice.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req userservice.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.JsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.Register(req)
	if err != nil {
		httputils.JsonError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		httputils.JsonError(w, err.Error(), http.StatusUnprocessableEntity)
	}
}
