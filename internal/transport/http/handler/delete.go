package handler

import (
	"net/http"

	"github.com/google/uuid"
)

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
