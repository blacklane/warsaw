# warsaw
Warsaw is a JSON Logger and context-based middleware to log HTTP requests (and more) in Golang projects (in [Kiev](https://github.com/blacklane/kiev) format).

## Development
 
The project is standalone [Go Module](https://blog.golang.org/using-go-modules). So everything as for regular gomodule applies, e.g. `go mod tidy`, `go build` and `go list -m all` commands. 
 
To run tests of the project use the `go test ./...` command. It will fetch all required dependencies and run all tests.

## Usage

There is variety of projects you might want to use this lib. 

* [HTTP requests](#http-requests)
* [Standalone app logger](#standalone-app-logger) (e.g. Kafka consumer project)
* [AWS Lambda](#aws-lambda)

But it all starts presumably with importing the dependency in your code with `import "github.com/blacklane/warsaw/logger"`. Yet it in details showcased in use-cases below and described in public [API reference](#api-reference) below.

### HTTP requests
Add the middleware to your `http.HandlerFunc` stack. Example:

```go
package main

import (
  "fmt"
  "net/http"

  "github.com/blacklane/warsaw/logger"
  "github.com/blacklane/warsaw/request_context"
)

func PingHandler(w http.ResponseWriter, req *http.Request) {
  log := logger.Get(req.Context())

  log.Event("ping_started").Str("some_field", "value").Int("some_int", 123).Send()
  fmt.Fprint(w, "ok")
}

func main() {
  loggerMiddleware := logger.NewKievRequestLogger("MyAppName")
  handlerWithContext := request_context.TrackerMiddleware(loggerMiddleware(PingHandler))
  routeHandler := request_context.RequestRouteTracker("ping", handlerWithContext)

  http.HandleFunc("/ping", routeHandler)
  http.ListenAndServe(":8080", nil)
}
```


Lines logged to STDOUT will have the following format:

```json
{
  "level": "info",
  "application": "MyAppName",
  "entry_point": true,
  "host": "localhost",
  "path": "/ping",
  "request_depth": 0,
  "request_id": "85bec38d-1057-4ec6-88af-bed508e98594",
  "route": "ping",
  "tree_path": "T",
  "verb": "GET",
  "timestamp": "2019-11-05T16:24:38+01:00",
  "event": "ping_started",
  "some_field": "value",
  "some_int": 123
}

{
  "level": "info",
  "application": "MyAppName",
  "entry_point": true,
  "host": "localhost",
  "path": "/ping",
  "request_depth": 0,
  "request_id": "85bec38d-1057-4ec6-88af-bed508e98594",
  "route": "ping",
  "tree_path": "T",
  "verb": "GET",
  "timestamp": "2019-11-05T16:24:38+01:00",
  "event": "request_finished",
  "body": "",
  "ip": "::1",
  "params": "",
  "user_agent": "curl/7.54.0",
  "status": 200,
  "request_duration": 0.366895
}
```

### Standalone app logger


### AWS Lambda

TODO Add explanation how to use it in lambda environment. 

## API Reference

Package `logger`:

* [logger.New(ctx context.Context, appName string) (Logger, context.Context)](#loggernewctx-contextcontext-appname-string-logger-contextcontext)
* [logger.NewStandalone(appName string) (Logger, context.Context)](#loggernewstandaloneappname-string-logger-contextcontext)
* [logger.DefaultLogger](#loggerdefaultlogger)
* [logger.LogSink](#loggerlogsink)
* [logger.Event(name string) *LoggedEvent](#loggereventname-string-loggedevent)
* [logger.Get(ctx context.Context) Logger](#loggergetctx-contextcontext-logger)

Interface `Logger` from package `logger`:
 
* [(Logger)Event(name string) *LoggedEvent](#loggereventname-string-loggedevent)
* [(Logger)WithScope(map[string]interface{})](#loggerwithscopemapstringinterface)

HTTP request handlers:

* [logger.NewKievRequestLogger(appName string) func(http.HandlerFunc) http.HandlerFunc](#loggernewkievrequestloggerappname-string-funchttphandlerfunc-httphandlerfunc) 
* [logger.LogErrorWithBody(ctx context.Context, err error, errName, responseBody string)](#loggerlogerrorwithbodyctx-contextcontext-err-error-errname-responsebody-string)

Package `request_context`:

* [request_context.TrackerMiddleware(next http.HandlerFunc) http.HandlerFunc](#request_contexttrackermiddlewarenext-httphandlerfunc-httphandlerfunc)
* [request_context.SetTrackerHeaders(ctx, &req.Header)](#request_contextsettrackerheadersctx-reqheader)
* [request_context.RequestRouteTracker(route string, next http.HandlerFunc) http.HandlerFunc](#request_contextrequestroutetrackerroute-string-next-httphandlerfunc-httphandlerfunc)

Sub-package `contexts` from `request_context`:

* [request_context/contexts.GetRequestID(ctx context.Context) string](#request_contextcontextsgetrequestidctx-contextcontext-string)
* [request_context/contexts.GetRequestRoute(ctx context.Context) string](#request_contextcontextsgetrequestroutectx-contextcontext-string)

--- 

### `logger.New(ctx context.Context, appName string) (Logger, context.Context)`

Creates a new logger and records it in the ctx. Plus sets the application name to appName value.

First returned value is the logger instance and second is enhanced context that includes the logger.

### `logger.NewStandalone(appName string) (Logger, context.Context)`

Creates a new logger based on `context.Background()` with provided appName for the logging context. Returns same values as plain `logger.New(...)` [function](#loggernewctx-contextcontext-appname-string-logger-contextcontext).

### `logger.DefaultLogger`

Default logger instance. It is used by `logger.Event(...)` function. It also has the [default](#loggerlogsink) `LogSink` set. Default logger cannot be overwritten and cannot have the scope updated.

### `logger.LogSink`

Allows to setup the logger output. By Default it's set to `os.Stdout`. But can be overwritten.     

### `logger.Event(name string) *LoggedEvent`

Logs message to the [DefaultLogger](#loggerdefaultlogger). Can be used directly to quickly log something to default LogSink. It only sets the `event`, `level` and `timestamp` fields. Plus anything set in the chain of methods.

Sample:

```go
logger.Event("myEventName").Str("some", "data").Int("statusCode", 1234).Send()
```

logs:

```json
{
  "level":"info",
  "timestamp":"2020-01-18T18:50:26.795165+01:00",
  "event":"myEventName",
  "some":"data",
  "statusCode": 1234 
}
```

### `logger.Get(ctx context.Context) Logger` 

Returns Logger instance from the argument representing current `context.Context`. Useful to get the logger instance downstream somewhere deep in your app. Then you pass just the context instance all way down and get the `Logger` out using this function.

### `(Logger)Event(name string) *LoggedEvent`

Logs single log line for an event with `name`. It can be called on result of `logger.Get()` method result. The API is same as for `zerolog.Event`. To persist the event you need to call `.Send()` on the returned value.

Sample:

```go
log := logger.Get(ctx) // to get the Logger instance from context.Context
log.Event("atlas_request").Int("response_status_code", resp.StatusCode).Dur("elapsed", time.Since(requestStart)).Str("url", fullUrl).Send()
```

### `(Logger)WithScope(map[string]interface{})`

Updates the context of all logged events with that Logger instance. 

Sample:

```go
log, _ := logger.NewStandalone("myApp")
log.WithScope(map[string]interface{}{"important": "yes", "code": 42})
log.Event("myEvent").Str("crucial", "sure").Send()
```

logs: 

```json
{
    "level":"info", 
    "application":"myApp",
    "code":42,
    "important":"yes",
    "timestamp":"2020-01-18T18:48:08.708609+01:00",
    "event":"myEvent", 
    "crucial":"sure"
}
```

So there is the shared scope from `WithScope` and anything defined inline.

### `logger.NewKievRequestLogger(appName string) func(http.HandlerFunc) http.HandlerFunc` 

Creates a middleware that can wrap `http.HandlerFunc` of your server with logger inside `request.Context()` 

Sample of usage in [HTTP requests use-case](#http-requests).

### `logger.LogErrorWithBody(ctx context.Context, err error, errName, responseBody string)`

When run will update logger context and any subsequent `Logger.Event(name)` will report the error.

⚠️ Keep in mind it won't write anything to the log! It will just remember in logger context that this particular error occurred.

The purpose is to log request error responses. Because such `request_finished` events will be then marked as not-successful.

Sample `http.HandlerFunc` with handling the error and logging it accordingly:

```go
func MyHandler(w http.ResponseWriter, req *http.Request) {
	params, err := buildParams(req)
	if err != nil {
		responseBody := err.Json()
		logger.LogErrorWithBody(req.Context(), err, "params.ValidationError", responseBody)
		w.WriteHeader(422)
		fmt.Fprint(w, responseBody)
		return
	}
	result := code.Run(req.Context(), params)

	fmt.Fprint(w, result.Json())
}

func main() {
    middleware := logger.NewKievRequestLogger("myApp")
    http.HandlerFunc(middleware(MyHandler))
    http.ListenAndServe(":12345", nil)
}
```

### `request_context.TrackerMiddleware(next http.HandlerFunc) http.HandlerFunc`

This is used inside of `logger.NewKievRequestLogger(...)` so if you use it to get middleware/handlerFunc wrapper you get it already. But if in your use-case for any reasons you just want to extract the Request context details like `RequestId`, `TreePath` and `RequestDepth` but without logging the call in kiev format you should use this TrackerMiddleware.

It will record the `RequestContext` instance in the context of the request + pass the RequestId to the response Headers accordingly. This makes sense in for example `PingHandler` use-case. When logging every ping/health-check is an overkill but you might be interested in having the correlation ID like `RequestId` available in the call in case of an issue.

### `request_context.SetTrackerHeaders(ctx, &req.Header)`

When you recorded the `RequestContext` with this method you can pass it to downstream net/http.Request calls. 

Sample:

```go
req, err := http.NewRequest("GET", fullUrl, postBody)
request_context.SetTrackerHeaders(ctx, &req.Header)

httpClient := &http.Client{}
resp, err := httpClient.Do(req) // the called fullUrl will be done with all correct headers passed from the originated request 
```

### `request_context.RequestRouteTracker(route string, next http.HandlerFunc) http.HandlerFunc`

To add extra context related to route name behind particular `http.HandlerFunc` you can wrap your call with this extra middleware and specify in first argument the name of the route. It will then be reported with every line of the logged events.  

Sample of use of `RequestRouteTracker` and `TrackerHandler` together but without the `NewKievRequestLogger`:

```go
func pingHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "ok")
}

func main() {
    route("/ping", pingHandler)
    http.ListenAndServe(":12345", nil)
}

func route(routeName string, handler http.HandlerFunc) {
	handlerWithContext := request_context.TrackerMiddleware(handler)
	routeHandler := request_context.RequestRouteTracker(routeName, handlerWithContext)
	http.HandleFunc(path, routeHandler)
}
```  

### `request_context/contexts.GetRequestID(ctx context.Context) string`

Use this method if you need to access the `RequestId` recorded using `TrackingMiddleware` directly in your code. 

### `request_context/contexts.GetRequestRoute(ctx context.Context) string`
 
 Use this method if you need to access the `RequestRoute` recorded using `RequestRouteTracker` directly in your code.