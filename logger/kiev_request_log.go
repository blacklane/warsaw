package logger

import (
	"net"
	"net/http"
	"time"

	"github.com/blacklane/warsaw/logger/kiev_fields"
)

const (
	xForwardedForHeader = "X-Forwarded-For"
)

func NewKievRequestLogger(appName string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, request *http.Request) {
			requestStart := time.Now()
			logger, loggingContext := NewRequestLogger(appName, request)

			ww := NewLoggedResponseWriter(w)
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
		}
	}
}

func ipAddress(request *http.Request) string {
	forwardedIp := request.Header.Get(xForwardedForHeader)
	if len(forwardedIp) == 0 {
		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			ip = request.RemoteAddr
		}
		return ip
	}
	return forwardedIp
}
