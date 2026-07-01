package http_test

import (
	"bytes"
	"context"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	service "github.com/example/subscriptions-service/internal/service/subscription"
	"github.com/example/subscriptions-service/internal/transport/http/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeService struct {
	createFn    func(context.Context, service.CreateCommand) (*domain.Subscription, error)
	getFn       func(context.Context, uuid.UUID) (*domain.Subscription, error)
	updateFn    func(context.Context, service.UpdateCommand) (*domain.Subscription, error)
	deleteFn    func(context.Context, uuid.UUID) error
	listFn      func(context.Context, service.ListFilter) ([]domain.Subscription, error)
	totalCostFn func(context.Context, service.CostFilter) (int, error)
}

func (f fakeService) Create(ctx context.Context, cmd service.CreateCommand) (*domain.Subscription, error) {
	if f.createFn == nil {
		return nil, domain.ErrNotFound
	}
	return f.createFn(ctx, cmd)
}
func (f fakeService) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	if f.getFn == nil {
		return nil, domain.ErrNotFound
	}
	return f.getFn(ctx, id)
}
func (f fakeService) Update(ctx context.Context, cmd service.UpdateCommand) (*domain.Subscription, error) {
	if f.updateFn == nil {
		return nil, domain.ErrNotFound
	}
	return f.updateFn(ctx, cmd)
}
func (f fakeService) Delete(ctx context.Context, id uuid.UUID) error {
	if f.deleteFn == nil {
		return domain.ErrNotFound
	}
	return f.deleteFn(ctx, id)
}
func (f fakeService) List(ctx context.Context, f2 service.ListFilter) ([]domain.Subscription, error) {
	if f.listFn == nil {
		return nil, nil
	}
	return f.listFn(ctx, f2)
}
func (f fakeService) TotalCost(ctx context.Context, f2 service.CostFilter) (int, error) {
	if f.totalCostFn == nil {
		return 0, nil
	}
	return f.totalCostFn(ctx, f2)
}

func newMux(svc fakeService) *nethttp.ServeMux {
	h := handler.NewSubscriptionHandler(svc, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	mux := nethttp.NewServeMux()
	h.Register(mux)
	return mux
}

func TestCreateSubscriptionBadJSON(t *testing.T) {
	mux := newMux(fakeService{})
	req := httptest.NewRequest(nethttp.MethodPost, "/api/v1/subscriptions", bytes.NewBufferString(`{"bad":`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusBadRequest, rec.Code)
}

func TestGetSubscriptionNotFound(t *testing.T) {
	svc := fakeService{
		getFn: func(_ context.Context, _ uuid.UUID) (*domain.Subscription, error) {
			return nil, domain.ErrNotFound
		},
	}
	mux := newMux(svc)
	req := httptest.NewRequest(nethttp.MethodGet, "/api/v1/subscriptions/"+uuid.NewString(), nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusNotFound, rec.Code)
}

func TestTotalCostInvalidPeriod(t *testing.T) {
	mux := newMux(fakeService{
		totalCostFn: func(_ context.Context, _ service.CostFilter) (int, error) {
			return 0, domain.ErrInvalidPeriod
		},
	})
	req := httptest.NewRequest(nethttp.MethodGet, "/api/v1/subscriptions/total-cost?from=12-2025&to=07-2025", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusBadRequest, rec.Code)
}

func TestRouterUnknownPath(t *testing.T) {
	mux := newMux(fakeService{})
	req := httptest.NewRequest(nethttp.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusNotFound, rec.Code)
}

func TestUpdateSubscriptionBadUUID(t *testing.T) {
	mux := newMux(fakeService{})
	req := httptest.NewRequest(nethttp.MethodPut, "/api/v1/subscriptions/bad", bytes.NewBufferString(`{"service_name":"Yandex Plus","price":400,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusBadRequest, rec.Code)
}

func TestUpdateSubscriptionNotFound(t *testing.T) {
	svc := fakeService{
		updateFn: func(context.Context, service.UpdateCommand) (*domain.Subscription, error) {
			return nil, domain.ErrNotFound
		},
	}
	mux := newMux(svc)
	req := httptest.NewRequest(nethttp.MethodPut, "/api/v1/subscriptions/"+uuid.NewString(), bytes.NewBufferString(`{"service_name":"Yandex Plus","price":400,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusNotFound, rec.Code)
}

func TestDeleteSubscriptionBadUUID(t *testing.T) {
	mux := newMux(fakeService{})
	req := httptest.NewRequest(nethttp.MethodDelete, "/api/v1/subscriptions/bad", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusBadRequest, rec.Code)
}

func TestDeleteSubscriptionNotFound(t *testing.T) {
	svc := fakeService{
		deleteFn: func(context.Context, uuid.UUID) error { return domain.ErrNotFound },
	}
	mux := newMux(svc)
	req := httptest.NewRequest(nethttp.MethodDelete, "/api/v1/subscriptions/"+uuid.NewString(), nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusNotFound, rec.Code)
}

func TestListSubscriptionsBadQuery(t *testing.T) {
	mux := newMux(fakeService{})
	req := httptest.NewRequest(nethttp.MethodGet, "/api/v1/subscriptions?limit=bad", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusBadRequest, rec.Code)
}

func TestListSubscriptionsError(t *testing.T) {
	svc := fakeService{
		listFn: func(context.Context, service.ListFilter) ([]domain.Subscription, error) {
			return nil, domain.ErrNotFound
		},
	}
	mux := newMux(svc)
	req := httptest.NewRequest(nethttp.MethodGet, "/api/v1/subscriptions", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusNotFound, rec.Code)
}
