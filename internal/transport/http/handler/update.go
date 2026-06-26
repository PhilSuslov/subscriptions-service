package handler

import (
	"net/http"

	"github.com/google/uuid"
)

func (h *SubscriptionHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var req subscriptionRequest
	if err := decodeSubscriptionRequest(r, &req); err != nil {
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
	writeSubscriptionResponse(w, http.StatusOK, toResponse(s))
}
