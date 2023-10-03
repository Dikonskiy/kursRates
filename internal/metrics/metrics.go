package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RequestCount    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	requestCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)

	return &Metrics{
		RequestCount:    requestCount,
		RequestDuration: requestDuration,
	}
}

func (m *Metrics) IncRequestCount(method, status string) {
	m.RequestCount.WithLabelValues(method, status).Inc()
}

func (m *Metrics) ObserveRequestDuration(method, status string, duration float64) {
	m.RequestDuration.WithLabelValues(method, status).Observe(duration)
}
