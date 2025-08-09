package controller

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) ReadSubscription(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	if idString == "" {
		sendError(w, "id parameter not set", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		sendError(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	subscription, err := h.aggregationService.ReadSubscription(id)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccess(w, subscription, http.StatusOK)
}
