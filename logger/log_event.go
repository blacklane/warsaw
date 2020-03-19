package logger

import (
	"github.com/blacklane/warsaw/logger/kiev_fields"
)

type eventLogger interface {
	Event(name string) *LoggedEvent
	ErrorEvent(name string, err error) *LoggedEvent
}

// Event logs single log line for an event with `name`. It can be called on result of `logger.Get()` method result.
// The API is same as for `zerolog.Event`. To persist the event you need to call `.Send()` on the returned value.
// Sample:
//      log := logger.Get(ctx)
//      log.Event("atlas_request").Int("response_status_code", resp.StatusCode).Dur("elapsed", time.Since(requestStart)).Str("url", fullUrl).Send()
func (logger logger) Event(name string) *LoggedEvent {
	return logger.log.Info().Timestamp().Str(kiev_fields.Event, name)
}

// ErrorEvent logs single log line for an event with the log level set to error. It passes the `err` parameter in to
// zerolog, which prints it as part of the log line.
func (logger logger) ErrorEvent(name string, err error) *LoggedEvent {
	return logger.log.Err(err).Timestamp().Str(kiev_fields.Event, name)
}

// Event package function logs message with DefaultLogger
func Event(name string) *LoggedEvent {
	return DefaultLogger.Event(name)
}

// Event package function logs message with DefaultLogger on loglevel error
func ErrorEvent(name string, err error) *LoggedEvent {
	return DefaultLogger.ErrorEvent(name, err)
}
