package health

import (
	"net/http"
)

// HealthHandler обрабатывает health check запросы
type HealthHandler struct{}

// NewHealthHandler создает новый health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Handle обрабатывает GET /health запрос
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"ok","service":"order"}`))
	if err != nil {
		// Логируем ошибку, но не можем вернуть её клиенту
		// так как заголовки уже отправлены
		_ = err // Игнорируем ошибку записи
	}
}
