package logger

import (
	"context"

	"github.com/blacklane/warsaw/logger/kiev_fields"
)

type logger struct {
	log internalLogger
}

type Logger interface {
	eventLogger
	WithScope(map[string]interface{})
	IsDisabled() bool
}

// Get returns Logger instance from the argument representing current `context.Context`. Useful to get the logger
// instance downstream somewhere deep in your app. Then you pass just the context instance all way down and get
// the `Logger` out using this function.
func Get(ctx context.Context) Logger {
	return logger{internalLoggerFromContext(ctx)}
}

// New creates a logger with appName specified and attaches it to the provided ctx and the enriched
// context is returned as second value. If already logger exists in the context it returns it as it is.
func New(ctx context.Context, appName string) (Logger, context.Context) {
	if existingLogger := Get(ctx); !existingLogger.IsDisabled() {
		return existingLogger, ctx
	}

	log := newInternalLogger(LogSink)
	loggingContext := log.WithContext(ctx)
	log.UpdateContext(func(c Context) Context {
		return c.Fields(map[string]interface{}{kiev_fields.Application: appName})
	})
	return logger{log}, loggingContext
}

// NewStandalone creates a logger with appName in fresh Context( `context.Background()` ) return as 2nd value
func NewStandalone(appName string) (Logger, context.Context) {
	return New(context.Background(), appName)
}

// WithScope allows to add new fields to existing logger context
func (logger logger) WithScope(fields map[string]interface{}) {
	logger.log.UpdateContext(func(c Context) Context {
		return c.Fields(fields)
	})
}

func (logger logger) IsDisabled() bool {
	return logger.log.GetLevel() == disabledLoggerLevel
}
