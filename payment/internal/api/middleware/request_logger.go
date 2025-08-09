package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

// RequestLogger логирует информацию о HTTP запросах
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем response writer для перехвата статуса
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Выполняем следующий обработчик
		next.ServeHTTP(ww, r)

		// Логируем информацию о запросе
		duration := time.Since(start)
		logger.Info(r.Context(), fmt.Sprintf("%s %s %d %v", r.Method, r.URL.Path, ww.statusCode, duration))
	})
}

// responseWriter обертка для http.ResponseWriter для перехвата статуса
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
