package controller

import (
	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"net/http"
)

// ListSubscriptions godoc
// @Summary Получить список подписок
// @Description Возвращает полный список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} entity.SubscriptionRequest "Список подписок"
// @Failure 500 {object} ErrorResponse "Ошибка получения данных"
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subs, err := h.aggregationService.ListSubscriptions(ctx)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subsResp := make([]*entity.SubscriptionRequest, 0, len(subs))
	for _, sub := range subs {
		subsResp = append(subsResp, entity.ParseSubscriptionToRequest(sub))
	}

	sendSuccess(w, subsResp, http.StatusOK)
}
