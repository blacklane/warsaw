package logger

import (
	"net"
	"net/http"
	"time"

	"github.com/blacklane/warsaw/logger/kiev_fields"
	"github.com/blacklane/warsaw/request_context"
)

const (
	xForwardedForHeader = "X-Forwarded-For"
)

// NewKievRequestLogger creates a middleware that can wrap `http.HandlerFunc` of your server with logger
// inside of the `request.Context()`.
func NewKievRequestLogger(appName string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return request_context.TrackerMiddleware(func(w http.ResponseWriter, request *http.Request) {
			requestStart := time.Now()
			logger, loggingContext := NewRequestLogger(appName, request)

			ww := newLoggedResponseWriter(w)
			defer func() {
				logger.
					Event(kiev_fields.RequestFinishedEvent).
					Fields(map[string]interface{}{
						kiev_fields.Params:    request.URL.RawQuery,
						kiev_fields.IP:        ipAddress(request),
						kiev_fields.UserAgent: request.UserAgent(),
						kiev_fields.Body:      ww.Body(),
					}).
					Int(kiev_fields.Status, ww.StatusCode()).
					Dur(kiev_fields.RequestDuration, time.Since(requestStart)).
					Send()
			}()

			next.ServeHTTP(ww, request.WithContext(loggingContext))
		})
	}
}

func ipAddress(request *http.Request) string {
	forwardedIP := request.Header.Get(xForwardedForHeader)
	if len(forwardedIP) == 0 {
		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			ip = request.RemoteAddr
		}
		return ip
	}
	return forwardedIP
}
