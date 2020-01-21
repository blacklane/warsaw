package logger

import (
	"context"
	"strings"
	"testing"
)

func Test_logger_WithScope(t *testing.T) {
	out := captureLogs(func() {
		log, _ := NewStandalone("appName")
		log.WithScope(map[string]interface{}{
			"someField":     "willBe Recorded",
			"anotherNumber": 432,
			"andBoolean":    true,
		})
		log.Event("myEventWithScope").Str("extra", "other value").Send()
		log.WithScope(map[string]interface{}{
			"more": "values",
		})
		log.Event("myEventWithScope").Bool("second", true).Send()
	})
	logLines := strings.SplitN(out, "\n", 2)

	matchEvent(t, logLines[0], map[string]string{"level": "info", "event": "myEventWithScope", "extra": "other value", "someField": "willBe Recorded", "anotherNumber": "432", "andBoolean": "true"})
	matchEvent(t, logLines[1], map[string]string{"level": "info", "event": "myEventWithScope", "second": "true", "more": "values", "someField": "willBe Recorded", "anotherNumber": "432", "andBoolean": "true"})
}

func helper(ctx context.Context, val string) {
	log := Get(ctx)
	log.Event("fromHelper").Str("data", val).Send()
}

func TestGet(t *testing.T) {
	t.Run("when no logger in context Get() returns disabledLogger which logs nothing", func(t *testing.T) {
		ctx := context.Background()

		out := captureLogs(func() {
			theLog := Get(ctx)
			theLog.Event("myEvent").Send()
			if !theLog.IsDisabled() {
				t.Errorf("not disabled logger from empty context!")
			}
		})
		if len(out) > 0 {
			t.Errorf("when there is no logger/disabled it shouldn't log anything")
		}
	})
}

func TestNewLogger(t *testing.T) {
	t.Run("Logger can be retrieved With Get() from context", func(t *testing.T) {
		ctx := context.Background()

		out := captureLogs(func() {
			log, loggingContext := New(ctx, "contextWiseLogger")
			log.WithScope(map[string]interface{}{"env": "global"})
			theLogger := Get(loggingContext)
			theLogger.Event("eventFromGet").Send()
			helper(loggingContext, "run1")
			helper(loggingContext, "run2")
			helper(ctx, "wouldn't Be Logged because uses context without logger")
		})

		logLines := strings.SplitN(out, "\n", 3)
		matchEvent(t, logLines[0], map[string]string{"level": "info", "event": "eventFromGet", "env": "global"})
		matchEvent(t, logLines[1], map[string]string{"level": "info", "event": "fromHelper", "env": "global", "data": "run1"})
		matchEvent(t, logLines[2], map[string]string{"level": "info", "event": "fromHelper", "env": "global", "data": "run2"})
	})

	t.Run("there can be only one Logger per context", func(t *testing.T) {
		ctx := context.Background()

		out := captureLogs(func() {
			_, ctxWithLogger := New(ctx, "contextWiseLogger")
			_, secondCtxWithLogger := New(ctxWithLogger, "contextWiseLogger2")
			theLogger := Get(ctxWithLogger)
			theLogger.Event("eventFromLoggingContext1").Send()
			otherLogger := Get(secondCtxWithLogger)
			otherLogger.Event("eventFromLoggingContext2").Send()
		})

		logLines := strings.SplitN(out, "\n", 2)
		matchEvent(t, logLines[0], map[string]string{"level": "info", "event": "eventFromLoggingContext1", "application": "contextWiseLogger"})
		matchEvent(t, logLines[1], map[string]string{"level": "info", "event": "eventFromLoggingContext2", "application": "contextWiseLogger"})
	})
}

func TestNewStandaloneLogger(t *testing.T) {
	t.Run("Logger with always fresh context", func(t *testing.T) {
		out := captureLogs(func() {
			log, loggingContext := NewStandalone("newStandalone1")
			log.WithScope(map[string]interface{}{"env": "dev"})
			theLogger := Get(loggingContext)
			helper(loggingContext, "firstLogger")
			log2, loggingContext2 := NewStandalone("newStandalone2")
			log2.Event("logMe").Msg("in")
			theLogger.Event("eventFromGet").Send()
			helper(loggingContext2, "secondLogger")
		})

		logLines := strings.SplitN(out, "\n", 4)
		matchEvent(t, logLines[0], map[string]string{"level": "info", "application": "newStandalone1", "event": "fromHelper", "env": "dev", "data": "firstLogger"})
		matchEvent(t, logLines[1], map[string]string{"level": "info", "application": "newStandalone2", "event": "logMe", "message": "in"})
		matchEvent(t, logLines[2], map[string]string{"level": "info", "application": "newStandalone1", "event": "eventFromGet", "env": "dev"})
		matchEvent(t, logLines[3], map[string]string{"level": "info", "application": "newStandalone2", "event": "fromHelper", "data": "secondLogger"})
	})
}
