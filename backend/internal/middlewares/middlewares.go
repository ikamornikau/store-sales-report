package middlewares

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
	"uuid"

	"github.com/ikamornikau/store-sales-report/backend/internal/logs"
	"github.com/ikamornikau/store-sales-report/backend/internal/metrics"
)

const requestIDHeader = "X-Request-ID"

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}

		w.Header().Set(requestIDHeader, requestID)

		ctx := logs.ContextWithRequestID(r.Context(), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithRequestDuration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		recorder := newStatusRecorder(w)

		defer func() {
			metrics.ObserveRequestDuration(r.Method, requestPath(r), recorder.status,
				time.Since(started),
			)
		}()

		next.ServeHTTP(recorder, r)
	})
}

func WithStatusCodes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := newStatusRecorder(w)
		defer func() {
			metrics.ObserveStatusCode(r.Method, requestPath(r), recorder.status)
		}()

		next.ServeHTTP(recorder, r)
	})
}

func WithRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := newStatusRecorder(w)

		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			slog.ErrorContext(r.Context(), "request panic recovered",
				slog.Any("panic", rec),
				slog.String("stack", string(debug.Stack())),
			)

			if !recorder.wroteHeader {
				http.Error(recorder,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(recorder, r)
	})
}

func requestPath(r *http.Request) string {
	if r.Pattern != "" {
		return r.Pattern
	}

	return r.URL.Path
}
