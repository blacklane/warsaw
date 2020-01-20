package logger

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestEvent(t *testing.T) {
	out := captureDefaultLogger(func() {
		Event("eventName").Str("title", "something").Int("code", 422).Bool("IsImportant", true).Dur("elapsed", time.Minute).Send()
	})

	matchEvent(t, out, map[string]string{"level": "info", "event": "eventName", "title": "something", "code": "422", "IsImportant": "true", "elapsed": "60000"})
}

func Test_logger_Event(t *testing.T) {
	out := captureLogs(func() {
		log, _ := NewStandalone("anApp")
		log.Event("myEventLogger").Msg("my msg")
	})

	matchEvent(t, out, map[string]string{"level": "info", "application": "anApp", "event": "myEventLogger", "message": "my msg"})
}

func parseEventOutput(out string) map[string]interface{} {
	var loggedEvent map[string]interface{}
	if err := json.Unmarshal([]byte(out), &loggedEvent); len(out) > 0 && err != nil {
		panic(err)
	}
	return loggedEvent
}

func compareEvent(actual map[string]interface{}, expected map[string]string) error {
	for k, v := range expected {
		if v != fmt.Sprintf("%v", actual[k]) {
			return fmt.Errorf("key `%s` should be: '%v' but is now '%v'", k, v, actual[k])
		}
	}

	if _, err := time.Parse(time.RFC3339Nano, actual["timestamp"].(string)); err != nil {
		return err
	}

	return nil
}

func matchEvent(t *testing.T, logLine string, expectedEvent map[string]string) {
	if err := compareEvent(parseEventOutput(logLine), expectedEvent); err != nil {
		t.Errorf("event log line is not equal to expected value- %v\n`%v`\n`%v`", err, logLine, expectedEvent)
	}
}
