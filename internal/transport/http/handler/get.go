package handler

import (
	"net/http"

	"github.com/google/uuid"
)

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
	
	writeSubscriptionResponse(w, http.StatusOK, toResponse(s))
}
