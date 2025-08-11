package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Ararat25/subscription-aggregation-service/internal/logger"
	"go.uber.org/zap"
)

// ErrorResponse описывает структуру ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"error message"` // строка для сообщения об ошибке
}

// StatusResponse описывает структуру ответа со статусом
type StatusResponse struct {
	Status string `json:"status" example:"success"` // статус ответа
}

// sendSuccess отправляет успешный JSON-ответ с указанным статусом
func sendSuccess(w http.ResponseWriter, data any, statusCode int) {
	if data != nil {
		respBytes, err := json.Marshal(data)
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(respBytes)
		if err != nil {
			logger.Log.Error("error writing response", zap.Error(err))
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(statusCode)
}

// sendError отправляет JSON-ответ с ошибкой
func sendError(w http.ResponseWriter, errMsg string, statusCode int) {
	w.WriteHeader(statusCode)

	resp := ErrorResponse{Error: errMsg}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		http.Error(w, errMsg, statusCode)
	}
}
