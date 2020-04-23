package middlewares

import (
	"bytes"
	"net"
	"net/http"
	"strings"

	"github.com/blacklane/warsaw/constants"
	"github.com/blacklane/warsaw/logger"
	"github.com/blacklane/warsaw/request_context/contexts"
)

// TODO(Anderson): add application name when creating the logger
func NewHttpHandlerLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := logger.Now()
			urlPath := strings.Split(r.URL.Path, "?")[0] // TODO: obfuscate query string values and show the keys
			ctx := r.Context()

			l := logger.Get(ctx)
			logFields := map[string]interface{}{
				constants.FieldEntryPoint:   isEntryPoint(r),
				constants.FieldRequestDepth: contexts.GetRequestDepth(ctx),
				constants.FieldRequestID:    contexts.GetRequestID(ctx),
				constants.FieldTreePath:     contexts.GetTreePath(ctx),
				constants.FieldRoute:        contexts.GetRequestRoute(ctx),
				constants.FieldParams:       r.URL.RawQuery,
				constants.FieldIP:           ipAddress(r),
				constants.FieldUserAgent:    r.UserAgent(),
				constants.FieldHost:         r.Host,
				constants.FieldVerb:         r.Method,
				constants.FieldPath:         r.URL.Path,
			}
			l.WithFields(logFields)

			ww := responseWriter{ResponseWriter: w, body: &bytes.Buffer{}}

			defer func() {
				l.WithFields(map[string]interface{}{
					constants.FieldRequestDuration: logger.Now().Sub(startTime),
					constants.FieldStatus:          ww.statusCode})
				l.Event(constants.EventRequestFinished).
					Msgf("%s %s", r.Method, urlPath)
			}()

			next.ServeHTTP(&ww, r)
		})
	}
}

func ipAddress(r *http.Request) string {
	forwardedIP := r.Header.Get(constants.HeaderForwardedFor)
	if len(forwardedIP) == 0 {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		return ip
	}
	return forwardedIP
}

func isEntryPoint(r *http.Request) bool {
	return len(r.Header.Get(constants.HeaderRequestID)) == 0
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w responseWriter) StatusCode() int {
	return w.statusCode
}

func (w responseWriter) Body() string {
	w.body.Reset()
	return w.body.String()
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *responseWriter) Write(buf []byte) (int, error) {
	if w.statusCode < http.StatusOK || w.statusCode > 299 {
		w.body.Write(buf)
	}
	return w.ResponseWriter.Write(buf)
}
