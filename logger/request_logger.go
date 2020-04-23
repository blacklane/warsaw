package logger

import (
	"context"
	"net"
	"net/http"

	"github.com/blacklane/warsaw/constants"
	"github.com/blacklane/warsaw/request_context/contexts"
)

// NewRequestLogger registers logger on the request and it's context and loads the logger scope from the request context.
func NewRequestLogger(appName string, req *http.Request) (Logger, context.Context) {
	log, loggingContext := New(req.Context(), appName)
	setupFromRequestContext(log, req)
	return log, loggingContext
}

func setupFromRequestContext(log Logger, req *http.Request) {
	ctx := req.Context()
	entryPoint := len(req.Header.Get(constants.HeaderRequestID)) == 0
	log.WithScope(map[string]interface{}{
		constants.FieldEntryPoint:   entryPoint,
		constants.FieldRequestID:    contexts.GetRequestID(ctx),
		constants.FieldRequestDepth: contexts.GetRequestDepth(ctx),
		constants.FieldTreePath:     contexts.GetTreePath(ctx),
		constants.FieldRoute:        contexts.GetRequestRoute(ctx),
		constants.FieldHost:         hostName(req),
		constants.FieldVerb:         req.Method,
		constants.FieldPath:         req.URL.Path,
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
	log := Get(ctx)
	log.WithScope(map[string]interface{}{
		constants.FieldErrorClass:   errName,
		constants.FieldErrorMessage: err,
		constants.FieldBody:         responseBody,
	})
}
