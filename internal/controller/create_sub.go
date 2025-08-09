package controller

import (
	"bytes"
	"encoding/json"
	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"net/http"
)

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
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

	id, err := h.aggregationService.CreateSubscription(newSub)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, struct {
		Id int64 `json:"id"`
	}{
		Id: id,
	}, http.StatusOK)
}
