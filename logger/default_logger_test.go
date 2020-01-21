package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
)

func Test_DefaultLogger(t *testing.T) {
	out := captureDefaultLogger(func() {
		DefaultLogger.Event("theEvent").Str("aTitle", "sth").Int("num", 422).Bool("valid", false).Err(fmt.Errorf("something failed")).Send()
	})

	matchEvent(t, out, map[string]string{"level": "info", "event": "theEvent", "aTitle": "sth", "num": "422", "valid": "false", "error": "something failed"})
}

func captureDefaultLogger(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = writer
	os.Stderr = writer
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	DefaultLogger = buildDefaultLogger(os.Stdout)
	f()
	writer.Close()
	return <-out
}
