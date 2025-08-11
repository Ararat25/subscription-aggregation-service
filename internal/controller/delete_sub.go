package controller

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// DeleteSubscription godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} StatusResponse "Статус выполнения"
// @Failure 400 {object} ErrorResponse "Некорректный ID или ошибка удаления"
// @Router /subscription/delete/{id} [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
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
	err = h.aggregationService.DeleteSubscription(ctx, id)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccess(w, StatusResponse{Status: "success"}, http.StatusOK)
}
