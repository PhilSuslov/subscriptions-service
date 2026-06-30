package handler

import (
	"encoding/json"
	"net/http"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
)

type subscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type subscriptionsListResponse struct {
	Items []subscriptionResponse `json:"items"`
}

type totalCostResponse struct {
	TotalCost int `json:"total_cost"`
}

type healthResponse struct {
	Status string `json:"status"`
}

func toResponse(s *domain.Subscription) subscriptionResponse {
	var end *string
	if s.EndMonth != nil {
		v := s.EndMonth.String()
		end = &v
	}

	return subscriptionResponse{
		ID:          s.ID.String(),
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID.String(),
		StartDate:   s.StartMonth.String(),
		EndDate:     end,
		CreatedAt:   s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func writeSubscriptionResponse(w http.ResponseWriter, status int, v subscriptionResponse) {
	b, _ := json.Marshal(v)
	writeJSON(w, status, b)
}

func writeSubscriptionsListResponse(w http.ResponseWriter, status int, v subscriptionsListResponse) {
	b, _ := json.Marshal(v)
	writeJSON(w, status, b)
}

func writeTotalCostResponse(w http.ResponseWriter, status int, v totalCostResponse) {
	b, _ := json.Marshal(v)
	writeJSON(w, status, b)
}

func writeHealthResponse(w http.ResponseWriter, status int) {
	b, _ := json.Marshal(healthResponse{Status: "ok"})
	writeJSON(w, status, b)
}

func writeError(w http.ResponseWriter, status int, err error) {
	b, _ := json.Marshal(errorResponse{Error: err.Error()})
	writeJSON(w, status, b)
}

func writeJSON(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}
