// Package http содержит HTTP-middleware: request id, логирование, recovery, метрики.
package http

import (
	"context"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/metrics"
)

// HeaderRequestID — заголовок, в котором ходит идентификатор запроса.
const HeaderRequestID = "X-Request-Id"

// statusRecorder запоминает код ответа: сам http.ResponseWriter его не отдаёт.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// RequestID кладёт идентификатор запроса в контекст и в заголовок ответа.
// Если клиент прислал свой — берём его, чтобы связать логи разных сервисов.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(HeaderRequestID)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		w.Header().Set(HeaderRequestID, requestID)

		ctx := context.WithValue(r.Context(), logger.TraceIDKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logging пишет в лог метод, путь, статус и длительность запроса.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		ctx := r.Context()

		defer func() {
			logger.Info(ctx, "http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rec.status),
				zap.Duration("duration", time.Since(start)),
			)
		}()

		next.ServeHTTP(rec, r)
	})
}

// Metrics считает количество и длительность HTTP-запросов.
// Путь берём из шаблона маршрута, иначе кардинальность метрики взорвётся от UUID.
func Metrics(routePattern func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rec, r)

			path := r.URL.Path
			if routePattern != nil {
				if pattern := routePattern(r); pattern != "" {
					path = pattern
				}
			}

			metrics.ObserveHTTP(r.Method, path, strconv.Itoa(rec.status), time.Since(start))
		})
	}
}

// Recoverer превращает панику в 500 и пишет стектрейс в лог.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		defer func() {
			if rvr := recover(); rvr != nil {
				logger.Error(ctx, "panic recovered",
					zap.Any("error", rvr),
					zap.String("stack", string(debug.Stack())),
				)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
