# Warsaw

[![pkg.go.dev](https://img.shields.io/badge/pkg.dev.go-reference-blue)](https://pkg.go.dev/github.com/blacklane/warsaw?tab=overview)

`warsaw` is a JSON Logger wrapper around [zerolog](https://github.com/rs/zerolog). Also a http middleware to log HTTP 
requests (and more) for Go projects.
The http request log follows the same format and fields as in [Kiev](https://github.com/blacklane/kiev), and you can add
more as you wish.

You can find the full api reference on [pkg.go.dev](https://pkg.go.dev/github.com/blacklane/warsaw?tab=overview) or
[GoDoc](https://godoc.org/github.com/blacklane/warsaw)

<!-- TOC depthFrom:1 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Install](#install)
- [Usage](#usage)
    - [HTTP requests](#http-requests)
    - [Standalone logger](#standalone-logger)
    - [AWS Lambda](#aws-lambda)
- [Development](#development)

<!-- /TOC -->

## Install

```go
go get -u github.com/blacklane/warsaw
```

## Usage

It log to the standard output (`os.Stdout`). All the outputs shown below are pretty printed for convenience, 
`warsaw`'s output is not pretty printed. All examples are in [examples](examples/)

### HTTP requests

Add the middleware to your `http.Handler`:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/blacklane/warsaw/logger"
)

func pingHandler(w http.ResponseWriter, req *http.Request) {
	log := logger.Get(req.Context())

	log.Event("ping_started").Str("some_field", "value").Int("some_int", 123).Send()

	_, _ = fmt.Fprint(w, "ping")
}

func addLogMiddleware() {
	loggerMiddleware := logger.NewKievRequestLogger("MyAppName")

	http.HandleFunc("/ping", loggerMiddleware(pingHandler))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
```

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

### Standalone logger

#### 1. Via `DefaultLogger`

```go
    logger.Event("myPlainEvent").Msg("log my line")
    logger.Event("myComplexEvent").Str("aString", "field").Int("num", 422).Bool("valid", false).Err(fmt.Errorf("something failed")).Send()
```

```json
{
  "level": "info",
  "timestamp": "2020-02-11T16:41:36.690388+01:00",
  "event": "mySimpleEvent",
  "message": "log my line"
}
{
  "level": "info",
  "timestamp": "2020-02-11T16:41:36.690518+01:00",
  "event": "myComplexEvent",
  "aString": "field",
  "num": 422,
  "valid": false,
  "error": "something failed",
  "message": "foo failed"
}
```  

#### 2. Standalone logger.

But you need to provide the application name and as a return you will get the logger instance + the context.Context
that contains the logger if you would like to pass it to some underlying functions.

```go
func standaloneLogger() {
	log, ctx := logger.New(context.TODO(), "myAppName")
	log.Event("An event").Msg("hello world")

	useLoggerFromContext(ctx)
}

func useLoggerFromContext(ctx context.Context) {
	log := logger.Get(ctx)
	log.Event("useLoggerFromContext").Bool("aTruth", true).Msg("hello :)")
}
```

Would log something like this:

```json
{
  "level": "info",
  "application": "myAppName",
  "timestamp": "2020-02-11T17:42:27.886908+01:00",
  "event": "An event",
  "message": "hello world"
}
{
  "level": "info",
  "application": "myAppName",
  "timestamp": "2020-02-11T17:55:42.20271+01:00",
  "event": "loggerFromContext",
  "aTruth": true,
  "message": "hello :)"
}
```

### AWS Lambda

For now there is no extra middleware for lambdaHandlers but it can be used with existing setup and simple instruction 
that builds the logger and registers it in the handler context.

```go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/blacklane/warsaw/logger"
)

type MyEvent struct {
	Name string `json:"name"`
}

func foo(ctx context.Context) {
	log := logger.Get(ctx)
	log.Event("I'm").Msg("insideMethod")
}

func handleRequest(ctx context.Context, name MyEvent) (string, error) {
	log, ctx := logger.NewLambdaLogger(ctx)
	log.Event("Called").Str("the_name", name.Name).Send()

	foo(ctx)
	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func main() {
	lambda.Start(handleRequest)
}
```

that would log events as follows:

```json
{
  "level": "info",
  "application": "myTestLambda",
  "entry_point": true,
  "lambda_function_arn": "arn:aws:lambda:eu-central-1:12345678901:function:myTestLambda",
  "lambda_function_version": "$LATEST",
  "lambda_memory_limit_in_mb": 512,
  "request_id": "4303499d-fccf-4cef-850c-26de69030463",
  "timestamp": "2020-01-20T14:47:05.491402536Z",
  "event": "Called",
  "the_name": "value1"
}
{
  "level": "info",
  "application": "myTestLambda",
  "entry_point": true,
  "lambda_function_arn": "arn:aws:lambda:eu-central-1:12345678901:function:myTestLambda",
  "lambda_function_version": "$LATEST",
  "lambda_memory_limit_in_mb": 512,
  "request_id": "4303499d-fccf-4cef-850c-26de69030463",
  "timestamp": "2020-01-20T14:47:05.491476597Z",
  "event": "I'm",
  "message": "insideMethod"
}
```


## Contributing

Open an issue, if the change is small, open the PR as well. If the change is big we commend to open the issue to discuss
it first. 

And don't forget:

```go
go test ./...
golint ./...
go vet ./...
```

## Licence

Copyright 2019 Blacklane

Licensed under [MIT License](LICENSE.md)