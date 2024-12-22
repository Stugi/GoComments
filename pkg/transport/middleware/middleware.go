package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "RequestID"

// Логирование
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()
		requestID, ok := ctx.Value(RequestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}

		log.Printf("[RequestID: %s] Started %s %s", requestID, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("[RequestID: %s] Completed in %v", requestID, time.Since(start))
	})
}

// RequestIDMiddleware добавляет сквозной идентификатор запроса в контекст
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String() // Генерируем новый ID
		}

		// Добавляем requestID в заголовок ответа
		w.Header().Set("X-Request-ID", requestID)

		// Сохраняем requestID в контекст
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
