package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	myError "github.com/Ararat25/subscription-aggregation-service/internal/error"
	"github.com/go-chi/chi/v5"
)

// ReadSubscription godoc
// @Summary Получить подписку по ID
// @Description Возвращает данные подписки по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} entity.SubscriptionRequest "Информация о подписке"
// @Failure 400 {object} ErrorResponse "Неверный параметр id или ошибка получения данных"
// @Router /subscription/{id} [get]
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

	ctx := r.Context()
	subscription, err := h.aggregationService.ReadSubscription(ctx, id)
	if errors.Is(err, myError.ErrSubscriptionNotFound) {
		sendError(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	subResp := entity.ParseSubscriptionToRequest(subscription)

	sendSuccess(w, subResp, http.StatusOK)
}
