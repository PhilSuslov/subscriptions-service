package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	app "github.com/example/subscriptions-service/internal/service/subscription"
	routerhttp "github.com/example/subscriptions-service/internal/transport/http"
	"github.com/example/subscriptions-service/internal/transport/http/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type smokeRepo struct {
	mu    sync.Mutex
	items map[uuid.UUID]*domain.Subscription
}

func newSmokeRepo() *smokeRepo {
	return &smokeRepo{items: make(map[uuid.UUID]*domain.Subscription)}
}

func (r *smokeRepo) Create(_ context.Context, s *domain.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *s
	r.items[s.ID] = &cp
	return nil
}

func (r *smokeRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *s
	return &cp, nil
}

func (r *smokeRepo) Update(_ context.Context, s *domain.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[s.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *s
	r.items[s.ID] = &cp
	return nil
}

func (r *smokeRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.items, id)
	return nil
}

func (r *smokeRepo) List(_ context.Context, f app.ListFilter) ([]domain.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]domain.Subscription, 0, len(r.items))
	for _, s := range r.items {
		if f.UserID != nil && s.UserID != *f.UserID {
			continue
		}
		if f.ServiceName != nil && *f.ServiceName != "" && !stringsEqualFold(*f.ServiceName, s.ServiceName) {
			continue
		}
		out = append(out, *s)
	}
	return out, nil
}

func (r *smokeRepo) TotalCost(_ context.Context, f app.CostFilter) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	total := 0
	for _, s := range r.items {
		if f.UserID != nil && s.UserID != *f.UserID {
			continue
		}
		if f.ServiceName != nil && *f.ServiceName != "" && !stringsEqualFold(*f.ServiceName, s.ServiceName) {
			continue
		}
		if s.StartMonth.After(f.PeriodTo) {
			continue
		}
		if s.EndMonth != nil && s.EndMonth.Before(f.PeriodFrom) {
			continue
		}
		total += s.Price
	}
	return total, nil
}

func stringsEqualFold(a, b string) bool {
	return strings.EqualFold(a, b)
}

func TestSmokeAllEndpoints(t *testing.T) {
	repo := newSmokeRepo()
	uc := app.NewUseCase(repo)
	h := handler.NewSubscriptionHandler(uc, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	srv := routerhttp.NewRouter(h, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	createBody := fmt.Sprintf(`{"service_name":"Yandex Plus","price":400,"user_id":"%s","start_date":"07-2025"}`, userID)

	createReq := httptest.NewRequest(nethttp.MethodPost, "/api/v1/subscriptions", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	srv.ServeHTTP(createRec, createReq)
	require.Equal(t, nethttp.StatusCreated, createRec.Code)

	var created struct {
		ID string `json:"id"`
	}
	require.NoError(t, json.NewDecoder(createRec.Body).Decode(&created))
	require.NotEmpty(t, created.ID)

	getReq := httptest.NewRequest(nethttp.MethodGet, "/api/v1/subscriptions/"+created.ID, nil)
	getRec := httptest.NewRecorder()
	srv.ServeHTTP(getRec, getReq)
	require.Equal(t, nethttp.StatusOK, getRec.Code)

	listReq := httptest.NewRequest(nethttp.MethodGet, "/api/v1/subscriptions?limit=10&offset=0", nil)
	listRec := httptest.NewRecorder()
	srv.ServeHTTP(listRec, listReq)
	require.Equal(t, nethttp.StatusOK, listRec.Code)

	updateBody := fmt.Sprintf(`{"service_name":"Yandex Music","price":500,"user_id":"%s","start_date":"07-2025"}`, userID)
	updateReq := httptest.NewRequest(nethttp.MethodPut, "/api/v1/subscriptions/"+created.ID, bytes.NewBufferString(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	srv.ServeHTTP(updateRec, updateReq)
	require.Equal(t, nethttp.StatusOK, updateRec.Code)

	totalReq := httptest.NewRequest(nethttp.MethodGet, "/api/v1/subscriptions/total-cost?from=07-2025&to=07-2025&user_id="+userID.String(), nil)
	totalRec := httptest.NewRecorder()
	srv.ServeHTTP(totalRec, totalReq)
	require.Equal(t, nethttp.StatusOK, totalRec.Code)
	require.Contains(t, totalRec.Body.String(), `"total_cost"`)

	deleteReq := httptest.NewRequest(nethttp.MethodDelete, "/api/v1/subscriptions/"+created.ID, nil)
	deleteRec := httptest.NewRecorder()
	srv.ServeHTTP(deleteRec, deleteReq)
	require.Equal(t, nethttp.StatusNoContent, deleteRec.Code)

	getAfterDeleteReq := httptest.NewRequest(nethttp.MethodGet, "/api/v1/subscriptions/"+created.ID, nil)
	getAfterDeleteRec := httptest.NewRecorder()
	srv.ServeHTTP(getAfterDeleteRec, getAfterDeleteReq)
	require.Equal(t, nethttp.StatusNotFound, getAfterDeleteRec.Code)

	healthReq := httptest.NewRequest(nethttp.MethodGet, "/health", nil)
	healthRec := httptest.NewRecorder()
	srv.ServeHTTP(healthRec, healthReq)
	require.Equal(t, nethttp.StatusOK, healthRec.Code)
}
