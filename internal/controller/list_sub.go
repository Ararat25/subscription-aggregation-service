package controller

import (
	"fmt"
	"net/http"
)

func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.aggregationService.ListSubscriptions()
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(subs)

	sendSuccess(w, subs, http.StatusOK)
}
