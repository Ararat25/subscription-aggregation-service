package controller

import (
	"bytes"
	"encoding/json"
	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"net/http"
)

// CreateSubscription godoc
// @Summary Создать новую подписку
// @Description Добавляет новую подписку в систему
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body entity.SubscriptionRequest true "Данные подписки"
// @Success 200 {object} map[string]int64 "ID созданной подписки"
// @Failure 400 {object} map[string]string "Некорректные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscription [post]
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

	ctx := r.Context()
	id, err := h.aggregationService.CreateSubscription(ctx, newSub)
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
