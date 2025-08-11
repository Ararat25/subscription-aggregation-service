package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/model"
	"github.com/google/uuid"
)

// TotalCostControllerResponse - структура для ответа от контроллера TotalCost
type TotalCostControllerResponse struct {
	TotalCost int `json:"total_cost" example:"2000"` // суммарная стоимость подписок
}

// TotalCost godoc
// @Summary Получить общую стоимость подписок
// @Description Возвращает суммарную стоимость подписок за указанный период с возможной фильтрацией по id пользователя и названию сервиса
// @Tags subscriptions
// @Produce json
// @Param from query string true "Дата начала периода (формат MM-YYYY)"
// @Param to query string true "Дата конца периода (формат MM-YYYY)"
// @Param id query string false "UUID пользователя"
// @Param service_name query string false "Название сервиса"
// @Success 200 {object} TotalCostControllerResponse "Общая стоимость"
// @Failure 400 {object} ErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/cost [get]
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

	ctx := r.Context()
	cost, err := h.aggregationService.TotalCost(ctx, from, to, id, serviceName)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, TotalCostControllerResponse{
		TotalCost: cost,
	}, http.StatusOK)
}
