package controller

import (
	"net/http"
)

// ListSubscriptions godoc
// @Summary Получить список подписок
// @Description Возвращает полный список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} entity.Subscription "Список подписок"
// @Failure 500 {object} ErrorResponse "Ошибка получения данных"
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subs, err := h.aggregationService.ListSubscriptions(ctx)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, subs, http.StatusOK)
}
