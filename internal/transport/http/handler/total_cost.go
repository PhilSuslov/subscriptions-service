package handler

import (
	"net/http"
	"strings"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	app "github.com/example/subscriptions-service/internal/service/subscription"
	"github.com/google/uuid"
)

func (h *SubscriptionHandler) totalCost(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from, err := domain.ParseMonth(q.Get("from"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	to, err := domain.ParseMonth(q.Get("to"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var userID *uuid.UUID
	if v := q.Get("user_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		userID = &id
	}
	var serviceName *string
	if v := strings.TrimSpace(q.Get("service_name")); v != "" {
		serviceName = &v
	}
	total, err := h.uc.TotalCost(r.Context(), app.CostFilter{UserID: userID, ServiceName: serviceName, PeriodFrom: from, PeriodTo: to})
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeTotalCostResponse(w, http.StatusOK, totalCostResponse{TotalCost: total})
}
