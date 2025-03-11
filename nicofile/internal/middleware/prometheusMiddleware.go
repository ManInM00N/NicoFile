package middleware

import (
	"net/http"
	"time"

	"main/nicofile/internal/metrics"
)

func PrometheusMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 包装 ResponseWriter 以捕获状态码
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		// 记录请求总数和响应时间
		metrics.RequestCount.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.status)).Inc()
		metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.status)).Observe(time.Since(start).Seconds())
	})
}

// 自定义 ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
