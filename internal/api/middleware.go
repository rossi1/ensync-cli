package api

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Middleware func(next http.RoundTripper) http.RoundTripper

type loggingTransport struct {
	next   http.RoundTripper
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &loggingTransport{
			next:   next,
			logger: logger,
		}
	}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.next.RoundTrip(req)

	duration := time.Since(start)

	t.logger.Info("API Request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Duration("duration", duration),
		zap.Int("status", resp.StatusCode),
	)

	return resp, err
}
