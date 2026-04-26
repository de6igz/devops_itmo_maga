package httpapi

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newMetrics() (*prometheus.Registry, echo.MiddlewareFunc) {
	registry := prometheus.NewRegistry()

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "game_catalog_http_requests_total",
			Help: "Total number of HTTP requests handled by the backend.",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "game_catalog_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	inFlightRequests := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "game_catalog_http_requests_in_flight",
			Help: "Current number of HTTP requests being handled by the backend.",
		},
	)

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		requestsTotal,
		requestDuration,
		inFlightRequests,
	)

	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/metrics" {
				return next(c)
			}

			startedAt := time.Now()
			inFlightRequests.Inc()
			defer inFlightRequests.Dec()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			status := c.Response().Status
			if status == 0 {
				status = 200
			}

			path := c.Path()
			if path == "" {
				path = c.Request().URL.Path
			}

			labels := prometheus.Labels{
				"method": c.Request().Method,
				"path":   path,
				"status": strconv.Itoa(status),
			}

			requestsTotal.With(labels).Inc()
			requestDuration.With(labels).Observe(time.Since(startedAt).Seconds())

			return nil
		}
	}

	return registry, middleware
}

func metricsHandler(registry *prometheus.Registry) echo.HandlerFunc {
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	return echo.WrapHandler(handler)
}
