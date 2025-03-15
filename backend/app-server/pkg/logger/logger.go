package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func New() error {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	return nil
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
	body   []byte
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	if r.status >= 400 { // Сохраняем тело только для ошибок
		r.body = make([]byte, len(b))
		copy(r.body, b)
	}
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
		Logger.Info("Request completed:",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.Int("status", lrw.status),
			zap.Duration("duration", duration),
		)

		if lrw.status >= 400 {
			Logger.Error("Request failed:",
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.Int("status", lrw.status),
				zap.Duration("duration", duration),
				zap.ByteString("response_body", lrw.body),
			)
		}
	})
}
