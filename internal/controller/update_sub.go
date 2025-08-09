package controller

import (
	"bytes"
	"encoding/json"
	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"net/http"
)

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newSub *entity.SubscriptionRequest
	err = json.Unmarshal(buf.Bytes(), &newSub)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.aggregationService.UpdateSubscription(newSub)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccess(w, StatusResponse{Status: "success"}, http.StatusOK)
}
