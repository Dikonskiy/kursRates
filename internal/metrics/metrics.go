package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RequestCount    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	SelectDuration  *prometheus.HistogramVec
	SelectCount     *prometheus.CounterVec
	InsertDuration  *prometheus.HistogramVec
	DeleteDuration  *prometheus.HistogramVec
	DeleteCount     *prometheus.CounterVec
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

	selectDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_select_duration_seconds",
			Help:    "Duration of database select queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	selectCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_select_count_total",
			Help: "Total number of database select queries",
		},
		[]string{"method", "status"},
	)

	insertDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_insert_duration_seconds",
			Help:    "Duration of database insert queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	deleteDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_delete_duration_seconds",
			Help:    "Duration of database delete queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	deleteCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_delete_count_total",
			Help: "Total number of database delete queries",
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(selectDuration)
	prometheus.MustRegister(selectCount)
	prometheus.MustRegister(insertDuration)
	prometheus.MustRegister(deleteDuration)
	prometheus.MustRegister(deleteCount)

	return &Metrics{
		RequestCount:    requestCount,
		RequestDuration: requestDuration,
		SelectDuration:  selectDuration,
		SelectCount:     selectCount,
		InsertDuration:  insertDuration,
		DeleteDuration:  deleteDuration,
		DeleteCount:     deleteCount,
	}
}

func (m *Metrics) IncRequestCount(method, status string) {
	m.RequestCount.WithLabelValues(method, status).Inc()
}

func (m *Metrics) ObserveRequestDuration(method, status string, duration float64) {
	m.RequestDuration.WithLabelValues(method, status).Observe(duration)
}

func (m *Metrics) ObserveSelectDuration(method, status string, duration float64) {
	m.SelectDuration.WithLabelValues(method, status).Observe(duration)
}

func (m *Metrics) IncSelectCount(method, status string) {
	m.SelectCount.WithLabelValues(method, status).Inc()
}

func (m *Metrics) ObserveInsertDuration(method, status string, duration float64) {
	m.InsertDuration.WithLabelValues(method, status).Observe(duration)
}

func (m *Metrics) IncDeleteCount(method, status string) {
	m.DeleteCount.WithLabelValues(method, status).Inc()
}

func (m *Metrics) ObserveDeleteDuration(method, status string, duration float64) {
	m.DeleteDuration.WithLabelValues(method, status).Observe(duration)
}
