package logger

import (
	"context"
	"net"
	"net/http"

	"github.com/blacklane/warsaw/logger/kiev_fields"
	"github.com/blacklane/warsaw/request_context/constants"
	"github.com/blacklane/warsaw/request_context/contexts"
)

func NewRequestLogger(appName string, req *http.Request) (Logger, context.Context) {
	log, loggingContext := New(req.Context(), appName)
	setupFromRequestContext(log, req)
	return log, loggingContext
}

func setupFromRequestContext(log Logger, req *http.Request) {
	ctx := req.Context()
	entryPoint := len(req.Header.Get(constants.RequestIDHeader)) == 0
	log.WithScope(map[string]interface{}{
		kiev_fields.EntryPoint:   entryPoint,
		kiev_fields.RequestID:    contexts.GetRequestID(ctx),
		kiev_fields.RequestDepth: contexts.GetRequestDepth(ctx),
		kiev_fields.TreePath:     contexts.GetTreePath(ctx),
		kiev_fields.Route:        contexts.GetRequestRoute(ctx),
		kiev_fields.Host:         hostName(req),
		kiev_fields.Verb:         req.Method,
		kiev_fields.Path:         req.URL.Path,
	})
}

func hostName(req *http.Request) string {
	host, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		host = req.Host
	}
	return host
}

// LogErrorWithBody updates logging context with details of the error and responseBody payload.
// This is later dumped with the request_finished event
func LogErrorWithBody(ctx context.Context, err error, errName, responseBody string) {
	internal, _ := Get(ctx).(logger)
	internal.log.UpdateContext(func(c Context) Context {
		return c.Str(kiev_fields.ErrorClass, errName).
			AnErr(kiev_fields.ErrorMessage, err).
			Str(kiev_fields.Body, responseBody)
	})
	// log := Get(ctx)
	// log.WithScope(map[string]interface{}{
	// 	kiev_fields.ErrorClass: errName,
	// 	kiev_fields.ErrorMessage: err,
	// 	kiev_fields.Body: responseBody,
	// })
}
