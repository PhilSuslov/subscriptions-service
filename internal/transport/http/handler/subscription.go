package handler

import (
	"log/slog"
	"net/http"
)

type SubscriptionHandler struct {
	uc  SubscriptionService
	log *slog.Logger
}

func NewSubscriptionHandler(uc SubscriptionService, log *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc, log: log}
}

func (h *SubscriptionHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("POST /api/v1/subscriptions", h.create)
	mux.HandleFunc("GET /api/v1/subscriptions", h.list)
	mux.HandleFunc("GET /api/v1/subscriptions/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/subscriptions/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/subscriptions/{id}", h.delete)
	mux.HandleFunc("GET /api/v1/subscriptions/total-cost", h.totalCost)
}
