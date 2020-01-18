package logger

// DefaultLogger represents everywhere available logger with default LogSink.
// It is also used with package logger.Event(...) method executions.
var DefaultLogger defaultLogger

type defaultLogger interface {
	Event(name string) *LoggedEvent
}

func init() {
	DefaultLogger = logger{log: newInternalLogger(LogSink)}
}

// Event package function logs message with DefaultLogger
func Event(name string) *LoggedEvent {
	return DefaultLogger.Event(name)
}
