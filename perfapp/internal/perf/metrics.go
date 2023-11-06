package perf

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTPRequestDuration HTTP request duration in seconds
var HTTPRequestDuration prometheus.Summary = promauto.NewSummary(prometheus.SummaryOpts{
	Name: "http_request_duration_seconds",
	Help: "Histogram of the HTTP requests duration",
})
