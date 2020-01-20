package logger

import (
	"fmt"
	"strings"
	"testing"
)

func TestLogErrorWithBody(t *testing.T) {
	t.Run("logs with extra context after recording error details", func(t *testing.T) {
		out := captureLogs(func() {
			log, ctx := NewStandalone("myApp")
			myErr := fmt.Errorf("My error")
			LogErrorWithBody(ctx, myErr, "errName", "some output")

			log.Event("myEvent").Send()
			log.Event("otherEvent").Msg("My content")
		})

		logLines := strings.SplitN(out, "\n", 2)
		matchEvent(t, logLines[0], map[string]string{"application": "myApp", "event": "myEvent", "error_message": "My error", "error_class": "errName", "body": "some output"})
		matchEvent(t, logLines[1], map[string]string{"application": "myApp", "event": "otherEvent", "error_message": "My error", "error_class": "errName", "body": "some output", "message": "My content"})
	})

}
