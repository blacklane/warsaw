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
`warsaw`'s output is not pretty printed. All examples are in [examples](examples/).

### HTTP requests

Add the middleware to a `http.Handler`:

```go
	loggerMiddleware := NewHttpHandlerLogger()

	h := loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = fmt.Fprint(w, "Hello, world") }))
	http.Handle("hello", h)

	// Output:
	// {"level":"info","application":"application name","entry_point":true,"host":"example.com","ip":"192.0.2.1","params":"","path":"/hello","request_depth":0,"request_id":"42","route":"","tree_path":"","user_agent":"","verb":"GET","request_duration":1000,"status":0,"timestamp":"2009-11-10T23:00:02Z","event":"request_finished","message":"GET /hello"}

```

```json
{
  "application": "application name",
  "entry_point": true,
  "event": "request_finished",
  "host": "example.com",
  "ip": "192.0.2.1",
  "level": "info",
  "message": "GET /hello",
  "params": "",
  "path": "/foo",
  "request_depth": 0,
  "request_duration": 1000,
  "request_id": "42",
  "route": "",
  "status": 0,
  "timestamp": "2009-11-10T23:00:02Z",
  "tree_path": "",
  "user_agent": "",
  "verb": "GET"
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