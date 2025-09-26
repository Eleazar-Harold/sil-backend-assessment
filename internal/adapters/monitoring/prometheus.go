package monitoring

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all the Prometheus metrics
type Metrics struct {
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge
	DatabaseConnections  prometheus.Gauge
	ActiveUsers          prometheus.Gauge
	CacheHits            *prometheus.CounterVec
	CacheMisses          *prometheus.CounterVec
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		HTTPRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		),
		DatabaseConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "database_connections_active",
				Help: "Number of active database connections",
			},
		),
		ActiveUsers: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users_total",
				Help: "Total number of active users",
			},
		),
		CacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache_type"},
		),
		CacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache_type"},
		),
	}
}

// PrometheusMiddleware creates a middleware for Prometheus metrics
func (m *Metrics) PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Increment in-flight requests
		m.HTTPRequestsInFlight.Inc()
		defer m.HTTPRequestsInFlight.Dec()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		m.HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		m.HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(wrapped.statusCode)).Inc()
	})
}

// GetHandler returns the Prometheus metrics handler
func (m *Metrics) GetHandler() http.Handler {
	return promhttp.Handler()
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// UpdateDatabaseConnections updates the database connections metric
func (m *Metrics) UpdateDatabaseConnections(count float64) {
	m.DatabaseConnections.Set(count)
}

// UpdateActiveUsers updates the active users metric
func (m *Metrics) UpdateActiveUsers(count float64) {
	m.ActiveUsers.Set(count)
}

// RecordCacheHit records a cache hit
func (m *Metrics) RecordCacheHit(cacheType string) {
	m.CacheHits.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss records a cache miss
func (m *Metrics) RecordCacheMiss(cacheType string) {
	m.CacheMisses.WithLabelValues(cacheType).Inc()
}
