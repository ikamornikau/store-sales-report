package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "code"})

	responses = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_responses_total",
		Help: "Number of HTTP responses by status code.",
	}, []string{"method", "path", "code"})
)

func ObserveRequestDuration(method, path string, statusCode int, duration time.Duration) {
	requestDuration.WithLabelValues(method, path, strconv.Itoa(statusCode)).Observe(duration.Seconds())
}

func ObserveStatusCode(method, path string, statusCode int) {
	responses.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
}
