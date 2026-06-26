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

func toResponse(s *domain.Subscription) subscriptionResponse {
	var end *string
	
	if s.EndMonth != nil {
		v := s.EndMonth.String()
		end = &v
	}

	return subscriptionResponse{ID: s.ID.String(), ServiceName: s.ServiceName, Price: s.Price, UserID: s.UserID.String(), StartDate: s.StartMonth.String(), EndDate: end, CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), UpdatedAt: s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

func writeSubscriptionResponse(w http.ResponseWriter, status int, v subscriptionResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeSubscriptionsListResponse(w http.ResponseWriter, status int, v subscriptionsListResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeTotalCostResponse(w http.ResponseWriter, status int, v totalCostResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeHealthResponse(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}
