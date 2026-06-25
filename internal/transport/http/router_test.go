package http

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/subscriptions-service/internal/transport/http/handler"
)

func TestNewRouter(t *testing.T) {
	h := handler.NewSubscriptionHandler(nil, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	router := NewRouter(h, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}
