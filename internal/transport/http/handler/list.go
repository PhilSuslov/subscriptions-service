package handler

import "net/http"

func (h *SubscriptionHandler) list(w http.ResponseWriter, r *http.Request) {
	f, err := parseListFilter(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	items, err := h.uc.List(r.Context(), f)
	if err != nil {
		h.handleError(w, err)
		return
	}

	resp := make([]subscriptionResponse, 0, len(items))
	for i := range items {
		resp = append(resp, toResponse(&items[i]))
	}

	writeSubscriptionsListResponse(w, http.StatusOK, subscriptionsListResponse{Items: resp})
}
