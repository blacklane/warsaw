# warsaw
JSON Logger and context-based middleware to log HTTP calls (and more) in Golang projects (in [Kiev](https://github.com/blacklane/kiev) format)

The purpose of this module is to provide HTTP request context-aware logger for Golang projects.

Usage requires adding some middlewares to your http.HandlerFunc stack. Sample:

```
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


Lines logged to STDOUT will be like:

```
{"level":"info","application":"MyAppName","entry_point":true,"host":"localhost","path":"/ping","request_depth":0,"request_id":"85bec38d-1057-4ec6-88af-bed508e98594","route":"ping","tree_path":"T","verb":"GET","timestamp":"2019-11-05T16:24:38+01:00","event":"ping_started","some_field":"value","some_int":123}
{"level":"info","application":"MyAppName","entry_point":true,"host":"localhost","path":"/ping","request_depth":0,"request_id":"85bec38d-1057-4ec6-88af-bed508e98594","route":"ping","tree_path":"T","verb":"GET","timestamp":"2019-11-05T16:24:38+01:00","event":"request_finished","body":"","ip":"::1","params":"","user_agent":"curl/7.54.0","status":200,"request_duration":0.366895}
```
