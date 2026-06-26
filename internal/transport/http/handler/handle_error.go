package handler

import (
	"errors"
	"net/http"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
)

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
