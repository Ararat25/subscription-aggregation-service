package controller

import (
	"github.com/Ararat25/subscription-aggregation-service/internal/model"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) TotalCost(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		sendError(w, "missing from or to parameter", http.StatusBadRequest)
		return
	}

	from, err := time.Parse(model.DateLayout, fromStr)
	if err != nil {
		sendError(w, "invalid from parameter", http.StatusBadRequest)
		return
	}
	to, err := time.Parse(model.DateLayout, toStr)
	if err != nil {
		sendError(w, "invalid to parameter", http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")
	serviceNameStr := r.URL.Query().Get("service_name")

	var serviceName *string
	if strings.TrimSpace(serviceNameStr) != "" {
		serviceName = &serviceNameStr
	}

	var id *uuid.UUID
	if idStr != "" {
		parsedID, err := uuid.Parse(idStr)
		if err != nil {
			sendError(w, "invalid id", http.StatusBadRequest)
			return
		}
		id = &parsedID
	}

	cost, err := h.aggregationService.TotalCost(from, to, id, serviceName)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, struct {
		TotalCost int
	}{
		TotalCost: cost,
	}, http.StatusOK)
}
