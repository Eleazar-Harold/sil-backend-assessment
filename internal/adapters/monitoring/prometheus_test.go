package monitoring

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// newTestMetrics creates a new metrics instance for testing with a custom registry
func newTestMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	factory := promauto.With(registry)

	return &Metrics{
		HTTPRequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		HTTPRequestsInFlight: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		),
		DatabaseConnections: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: "database_connections_active",
				Help: "Number of active database connections",
			},
		),
		ActiveUsers: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users_total",
				Help: "Total number of active users",
			},
		),
		CacheHits: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache_type"},
		),
		CacheMisses: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache_type"},
		),
	}
}

func TestNewMetrics(t *testing.T) {
	metrics := newTestMetrics()

	if metrics == nil {
		t.Error("Expected metrics to be non-nil")
	}

	// Check that all metrics are initialized
	if metrics.HTTPRequestsTotal == nil {
		t.Error("Expected HTTPRequestsTotal to be initialized")
	}
	if metrics.HTTPRequestDuration == nil {
		t.Error("Expected HTTPRequestDuration to be initialized")
	}
	if metrics.HTTPRequestsInFlight == nil {
		t.Error("Expected HTTPRequestsInFlight to be initialized")
	}
	if metrics.DatabaseConnections == nil {
		t.Error("Expected DatabaseConnections to be initialized")
	}
	if metrics.ActiveUsers == nil {
		t.Error("Expected ActiveUsers to be initialized")
	}
	if metrics.CacheHits == nil {
		t.Error("Expected CacheHits to be initialized")
	}
	if metrics.CacheMisses == nil {
		t.Error("Expected CacheMisses to be initialized")
	}
}

func TestMetrics_PrometheusMiddleware(t *testing.T) {
	metrics := newTestMetrics()

	t.Run("Middleware processes request", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test response"))
		})

		middleware := metrics.PrometheusMiddleware(handler)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got: %d", http.StatusOK, w.Code)
		}
		if w.Body.String() != "test response" {
			t.Errorf("Expected 'test response', got: %s", w.Body.String())
		}
	})

	t.Run("Middleware handles different HTTP methods", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		})

		middleware := metrics.PrometheusMiddleware(handler)

		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
		for _, method := range methods {
			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				t.Errorf("Expected status %d for %s, got: %d", http.StatusCreated, method, w.Code)
			}
		}
	})

	t.Run("Middleware handles different status codes", func(t *testing.T) {
		statusCodes := []int{200, 201, 400, 401, 404, 500}

		for _, statusCode := range statusCodes {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(statusCode)
			})

			middleware := metrics.PrometheusMiddleware(handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			if w.Code != statusCode {
				t.Errorf("Expected status %d, got: %d", statusCode, w.Code)
			}
		}
	})

	t.Run("Middleware handles panics gracefully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		middleware := metrics.PrometheusMiddleware(handler)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		// The middleware should handle panics gracefully due to the defer
		defer func() {
			if r := recover(); r != nil {
				// This is expected - the middleware doesn't handle panics
				t.Logf("Panic recovered as expected: %v", r)
			}
		}()

		// This will panic, but that's expected behavior
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic to occur")
				}
			}()
			middleware.ServeHTTP(w, req)
		}()
	})
}

func TestMetrics_GetHandler(t *testing.T) {
	metrics := newTestMetrics()
	handler := metrics.GetHandler()

	if handler == nil {
		t.Error("Expected handler to be non-nil")
	}

	// Test that the handler can be called
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// The handler should return some metrics data
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, w.Code)
	}
}

func TestResponseWriter(t *testing.T) {
	t.Run("WriteHeader captures status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Test default status code
		if rw.statusCode != http.StatusOK {
			t.Errorf("Expected default status code %d, got: %d", http.StatusOK, rw.statusCode)
		}

		// Test setting status code
		expectedStatus := http.StatusNotFound
		rw.WriteHeader(expectedStatus)

		if rw.statusCode != expectedStatus {
			t.Errorf("Expected status code %d, got: %d", expectedStatus, rw.statusCode)
		}
	})

	t.Run("WriteHeader calls underlying WriteHeader", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		expectedStatus := http.StatusInternalServerError
		rw.WriteHeader(expectedStatus)

		// Check that the underlying ResponseWriter also received the status code
		if w.Code != expectedStatus {
			t.Errorf("Expected underlying writer status code %d, got: %d", expectedStatus, w.Code)
		}
	})
}

func TestMetrics_UpdateDatabaseConnections(t *testing.T) {
	metrics := newTestMetrics()

	t.Run("Update database connections", func(t *testing.T) {
		expectedCount := 5.0
		metrics.UpdateDatabaseConnections(expectedCount)

		// We can't easily test the actual metric value without exposing it,
		// but we can test that the method doesn't panic
		metrics.UpdateDatabaseConnections(10.0)
		metrics.UpdateDatabaseConnections(0.0)
		metrics.UpdateDatabaseConnections(-1.0)
	})
}

func TestMetrics_UpdateActiveUsers(t *testing.T) {
	metrics := newTestMetrics()

	t.Run("Update active users", func(t *testing.T) {
		expectedCount := 100.0
		metrics.UpdateActiveUsers(expectedCount)

		// We can't easily test the actual metric value without exposing it,
		// but we can test that the method doesn't panic
		metrics.UpdateActiveUsers(50.0)
		metrics.UpdateActiveUsers(0.0)
		metrics.UpdateActiveUsers(1000.0)
	})
}

func TestMetrics_RecordCacheHit(t *testing.T) {
	metrics := newTestMetrics()

	t.Run("Record cache hit", func(t *testing.T) {
		cacheTypes := []string{"redis", "memory", "database"}

		for _, cacheType := range cacheTypes {
			metrics.RecordCacheHit(cacheType)
			// Record multiple hits
			metrics.RecordCacheHit(cacheType)
			metrics.RecordCacheHit(cacheType)
		}
	})

	t.Run("Record cache hit with empty cache type", func(t *testing.T) {
		metrics.RecordCacheHit("")
	})
}

func TestMetrics_RecordCacheMiss(t *testing.T) {
	metrics := newTestMetrics()

	t.Run("Record cache miss", func(t *testing.T) {
		cacheTypes := []string{"redis", "memory", "database"}

		for _, cacheType := range cacheTypes {
			metrics.RecordCacheMiss(cacheType)
			// Record multiple misses
			metrics.RecordCacheMiss(cacheType)
			metrics.RecordCacheMiss(cacheType)
		}
	})

	t.Run("Record cache miss with empty cache type", func(t *testing.T) {
		metrics.RecordCacheMiss("")
	})
}

func TestMetrics_Integration(t *testing.T) {
	metrics := newTestMetrics()

	t.Run("Full workflow", func(t *testing.T) {
		// Update various metrics
		metrics.UpdateDatabaseConnections(5.0)
		metrics.UpdateActiveUsers(100.0)
		metrics.RecordCacheHit("redis")
		metrics.RecordCacheMiss("memory")

		// Test middleware with a handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond) // Simulate some processing time
			w.WriteHeader(http.StatusOK)
		})

		middleware := metrics.PrometheusMiddleware(handler)

		req := httptest.NewRequest("POST", "/api/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got: %d", http.StatusOK, w.Code)
		}

		// Test metrics handler
		metricsHandler := metrics.GetHandler()
		req = httptest.NewRequest("GET", "/metrics", nil)
		w = httptest.NewRecorder()

		metricsHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected metrics handler status %d, got: %d", http.StatusOK, w.Code)
		}
	})
}

func TestMetrics_ConcurrentAccess(t *testing.T) {
	metrics := newTestMetrics()

	t.Run("Concurrent metric updates", func(t *testing.T) {
		done := make(chan bool, 10)

		// Start multiple goroutines updating metrics concurrently
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()

				metrics.UpdateDatabaseConnections(float64(i))
				metrics.UpdateActiveUsers(float64(i * 10))
				metrics.RecordCacheHit("redis")
				metrics.RecordCacheMiss("memory")
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}
