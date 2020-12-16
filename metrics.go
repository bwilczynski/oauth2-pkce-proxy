package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const metricsNamespace = "oauth2_pkce_proxy"

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Duration of HTTP requests by path and status code.",
	}, []string{"code", "method", "path"})
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricsNamespace,
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests by path and status code.",
	}, []string{"code", "method", "path"})
)

func instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wi := &responseWriterInterceptor{
			statusCode:     http.StatusOK,
			ResponseWriter: w,
		}

		start := time.Now()

		defer func() {
			duration := time.Since(start).Seconds()
			code := strconv.Itoa(wi.statusCode)
			method := strings.ToUpper(r.Method)
			path := r.URL.Path

			httpDuration.WithLabelValues(code, method, path).Observe(duration)
			httpRequestsTotal.WithLabelValues(code, method, path).Inc()
		}()

		next.ServeHTTP(wi, r)
	})
}

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
