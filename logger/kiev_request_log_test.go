package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blacklane/warsaw/request_context/constants"
)

func TestSoleNewKievRequestLogger(t *testing.T) {
	expectedReqId := "request-UUID"
	req, err := http.NewRequest("GET", "/ping", nil)
	req.Header.Set(constants.RequestIDHeader, expectedReqId)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	noop := func(w http.ResponseWriter, request *http.Request) {}

	out := captureLogs(func() {
		handler := NewKievRequestLogger("testApp")(noop)
		handler.ServeHTTP(rr, req)
	})

	if reqId := rr.Header().Get(constants.RequestIDHeader); reqId != expectedReqId {
		t.Errorf("handler returned unexpected X-Request-Id header: got `%v` want %v", reqId, expectedReqId)
	}

	var loggedEvent map[string]interface{}
	if err := json.Unmarshal([]byte(out), &loggedEvent); len(out) > 0 && err != nil {
		panic(err)
	}
	expectedEvent := map[string]string{"level": "info", "application": "testApp", "entry_point": "false", "path": "/ping", "request_depth": "0", "request_id": "request-UUID", "route": "", "tree_path": "T", "verb": "GET", "event": "request_finished", "body": "", "ip": "", "params": "", "user_agent": "", "status": "200"}
	if err := compareEvent(loggedEvent, expectedEvent); err != nil {
		t.Errorf("events are not equal - %v\n`%v`\n`%v`", err, loggedEvent, expectedEvent)
	}
	print(out)
}

func captureLogs(f func()) string {
	oldSink := LogSink
	var buf bytes.Buffer
	LogSink = &buf
	f()
	LogSink = oldSink

	return buf.String()
}

func compareEvent(actual map[string]interface{}, expected map[string]string) error {
	fmt.Println(actual, expected)
	for k, v := range expected {
		if v != fmt.Sprintf("%v", actual[k]) {
			return fmt.Errorf("key `%s` should be: '%v' but is now '%v'", k, v, actual[k])
		}
	}

	if _, err := time.Parse(time.RFC3339Nano, actual["timestamp"].(string)); err != nil {
		return err
	}

	if _, ok := actual["request_duration"].(float64); !ok {
		return fmt.Errorf("could not parse request_duration as float: %v", actual["request_duration"])
	}
	return nil
}
