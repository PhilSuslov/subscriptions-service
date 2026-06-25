package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	app "github.com/example/subscriptions-service/internal/service/subscription"
	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	uc  *app.UseCase
	log *slog.Logger
}

func NewSubscriptionHandler(uc *app.UseCase, log *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc, log: log}
}

type subscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

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

func (h *SubscriptionHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("POST /api/v1/subscriptions", h.create)
	mux.HandleFunc("GET /api/v1/subscriptions", h.list)
	mux.HandleFunc("GET /api/v1/subscriptions/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/subscriptions/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/subscriptions/{id}", h.delete)
	mux.HandleFunc("GET /api/v1/subscriptions/total-cost", h.totalCost)
}

func (h *SubscriptionHandler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *SubscriptionHandler) create(w http.ResponseWriter, r *http.Request) {
	var req subscriptionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	cmd, err := toCreateCommand(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s, err := h.uc.Create(r.Context(), cmd)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(s))
}

func (h *SubscriptionHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s, err := h.uc.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(s))
}

func (h *SubscriptionHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var req subscriptionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	cmd, err := toUpdateCommand(id, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s, err := h.uc.Update(r.Context(), cmd)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(s))
}

func (h *SubscriptionHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.uc.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) list(w http.ResponseWriter, r *http.Request) {
	f, err := parseListFilter(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	items, err := h.uc.List(r.Context(), f)
	if err != nil {
		h.handleError(w, err)
		return
	}
	resp := make([]subscriptionResponse, 0, len(items))
	for i := range items {
		resp = append(resp, toResponse(&items[i]))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": resp})
}

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
	writeJSON(w, http.StatusOK, map[string]int{"total_cost": total})
}

func parseListFilter(r *http.Request) (app.ListFilter, error) {
	q := r.URL.Query()
	f := app.ListFilter{Limit: 50}
	if v := q.Get("user_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			return f, err
		}
		f.UserID = &id
	}
	if v := strings.TrimSpace(q.Get("service_name")); v != "" {
		f.ServiceName = &v
	}
	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return f, err
		}
		f.Limit = n
	}
	if v := q.Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return f, err
		}
		f.Offset = n
	}
	return f, nil
}

func toCreateCommand(req subscriptionRequest) (app.CreateCommand, error) {
	id, err := uuid.Parse(req.UserID)
	if err != nil {
		return app.CreateCommand{}, err
	}
	return app.CreateCommand{ServiceName: req.ServiceName, Price: req.Price, UserID: id, StartDate: req.StartDate, EndDate: req.EndDate}, nil
}

func toUpdateCommand(id uuid.UUID, req subscriptionRequest) (app.UpdateCommand, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return app.UpdateCommand{}, err
	}
	return app.UpdateCommand{ID: id, ServiceName: req.ServiceName, Price: req.Price, UserID: userID, StartDate: req.StartDate, EndDate: req.EndDate}, nil
}

func toResponse(s *domain.Subscription) subscriptionResponse {
	var end *string
	if s.EndMonth != nil {
		v := s.EndMonth.String()
		end = &v
	}
	return subscriptionResponse{ID: s.ID.String(), ServiceName: s.ServiceName, Price: s.Price, UserID: s.UserID.String(), StartDate: s.StartMonth.String(), EndDate: end, CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), UpdatedAt: s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("request body must contain a single JSON object")
		}
		return fmt.Errorf("request body must contain a single JSON object")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, errorResponse{Error: err.Error()})
}

func (h *SubscriptionHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, domain.ErrInvalidPrice), errors.Is(err, domain.ErrInvalidServiceName), errors.Is(err, domain.ErrInvalidPeriod):
		writeError(w, http.StatusBadRequest, err)
	default:
		h.log.Error("request failed", "error", err)
		writeError(w, http.StatusInternalServerError, errors.New("internal error"))
	}
}
