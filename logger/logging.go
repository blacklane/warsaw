package logger

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	"github.com/blacklane/warsaw/logger/kiev_fields"
	"github.com/blacklane/warsaw/request_context/constants"
	"github.com/blacklane/warsaw/request_context/contexts"
)

type internalLogger interface {
	UpdateContext(func(zerolog.Context) zerolog.Context)
	Info() *zerolog.Event
}

type logger struct {
	log internalLogger
}

type Event = zerolog.Event

func (logger logger) Event(name string) *Event {
	return logger.log.Info().Timestamp().Str(kiev_fields.Event, name)
}

func Get(ctx context.Context) logger {
	return logger{zerolog.Ctx(ctx)}
}

func NewRequestLogger(appName string, req *http.Request) (logger, context.Context) {
	zerolog.TimestampFieldName = kiev_fields.Timestamp
	log := zerolog.New(os.Stdout)
	loggingContext := log.WithContext(req.Context())
	setRequestContext(&log, appName, req)
	return logger{log: &log}, loggingContext
}

func setRequestContext(log internalLogger, appName string, req *http.Request) {
	log.UpdateContext(func(c zerolog.Context) zerolog.Context {
		ctx := req.Context()
		entryPoint := len(req.Header.Get(constants.RequestIDHeader)) == 0

		return c.Fields(map[string]interface{}{
			kiev_fields.Application:  appName,
			kiev_fields.EntryPoint:   entryPoint,
			kiev_fields.RequestID:    contexts.GetRequestID(ctx),
			kiev_fields.RequestDepth: contexts.GetRequestDepth(ctx),
			kiev_fields.TreePath:     contexts.GetTreePath(ctx),
			kiev_fields.Route:        contexts.GetRequestRoute(ctx),
			kiev_fields.Host:         hostName(req),
			kiev_fields.Verb:         req.Method,
			kiev_fields.Path:         req.URL.Path,
		})

	})
}

func hostName(req *http.Request) string {
	host, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		host = req.Host
	}
	return host
}

func LogErrorWithBody(ctx context.Context, err error, errName, responseBody string) {
	Get(ctx).log.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(kiev_fields.ErrorClass, errName).
			AnErr(kiev_fields.ErrorMessage, err).
			Str(kiev_fields.Body, responseBody)
	})
}
