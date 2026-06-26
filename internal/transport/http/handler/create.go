package handler

import "net/http"

func (h *SubscriptionHandler) create(w http.ResponseWriter, r *http.Request) {
	var req subscriptionRequest
	if err := decodeSubscriptionRequest(r, &req); err != nil {
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
	writeSubscriptionResponse(w, http.StatusCreated, toResponse(s))
}
