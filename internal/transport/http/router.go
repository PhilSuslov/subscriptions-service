package http

import (
	"log/slog"
	nethttp "net/http"

	"github.com/example/subscriptions-service/internal/transport/http/handler"
	"github.com/example/subscriptions-service/internal/transport/http/middleware"
)

func NewRouter(h *handler.SubscriptionHandler, log *slog.Logger) nethttp.Handler {
	mux := nethttp.NewServeMux()
	h.Register(mux)
	return middleware.Recover(log)(middleware.Logging(log)(mux))
}
