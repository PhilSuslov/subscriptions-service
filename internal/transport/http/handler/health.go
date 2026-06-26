package handler

import "net/http"

func (h *SubscriptionHandler) health(w http.ResponseWriter, _ *http.Request) {
	writeHealthResponse(w, http.StatusOK)
}
