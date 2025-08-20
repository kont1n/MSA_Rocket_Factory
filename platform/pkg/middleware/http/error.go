package http

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse структура для возврата ошибок в JSON формате
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail детали ошибки
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeErrorResponse записывает структурированный JSON ответ с ошибкой
func writeErrorResponse(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Если не удалось записать JSON ответ, записываем простой текст
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error")) //nolint:errcheck,gosec
	}
}
