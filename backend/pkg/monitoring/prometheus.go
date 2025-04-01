package monitoring

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// Fiber için özel Prometheus middleware
func PrometheusMiddleware() fiber.Handler {
	reqCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Toplam HTTP istek sayısı",
		},
		[]string{"method", "path"},
	)

	reqDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP istek süreleri",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	prometheus.MustRegister(reqCount, reqDuration)

	return func(c *fiber.Ctx) error {
		path := c.Path()
		method := c.Method()

		timer := prometheus.NewTimer(reqDuration.WithLabelValues(method, path))
		defer timer.ObserveDuration()

		err := c.Next()
		reqCount.WithLabelValues(method, path).Inc()
		return err
	}
}

// /metrics endpoint'ini Fiber ile kullanmak için handler
func MetricsHandler() fiber.Handler {
	h := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())

	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/plain")
		h(c.Context()) // Fiber'ın context'ini Prometheus handler'ına yönlendir
		return nil
	}
}
