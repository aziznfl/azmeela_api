package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HttpRequestsTotal tracks total HTTP requests with labels for method, path, and status code
	HttpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "azmeela_http_requests_total",
		Help: "Total number of HTTP requests processed, partitioned by status code, method and path.",
	}, []string{"method", "path", "status"})

	// HttpRequestDuration tracks the latency of HTTP requests
	HttpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "azmeela_http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)
