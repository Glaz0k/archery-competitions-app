package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Logger, _ = zap.NewProduction()

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}

func LogMiddleware(hand http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := loggingResponseWriter{
			ResponseWriter: w,
			size:           0,
			status:         200,
		}
		hand.ServeHTTP(&lrw, r)
		duration := time.Since(start)
		Logger.Info("Request completed:", zap.String("method", r.Method),
			zap.String("url", r.URL.Path), zap.Int("status", lrw.status), zap.Duration("duration", duration))
	})
}
