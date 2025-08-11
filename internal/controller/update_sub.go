package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	myError "github.com/Ararat25/subscription-aggregation-service/internal/error"
)

// UpdateSubscription godoc
// @Summary Обновить подписку
// @Description Обновляет данные существующей подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body entity.SubscriptionRequest true "Данные подписки"
// @Success 200 {object} StatusResponse "Успешное обновление"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscription/update [put]
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

	err = validate.Struct(newSub)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.aggregationService.UpdateSubscription(ctx, newSub)
	if errors.Is(err, myError.ErrSubscriptionNotFound) {
		sendError(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccess(w, StatusResponse{Status: "success"}, http.StatusOK)
}
