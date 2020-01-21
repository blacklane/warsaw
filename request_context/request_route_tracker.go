package request_context

import (
	"net/http"

	"github.com/blacklane/warsaw/request_context/contexts"
)

// RequestRouteTracker is an alias for RouteTracker function.
//
// Deprecated: still used for historical compatibility. Please Use the RouteTracker() instead.
func RequestRouteTracker(route string, next http.HandlerFunc) http.HandlerFunc {
	return RouteTracker(route, next)
}

// RouteTracker stores the routeName in the logger context for future lookups
func RouteTracker(route string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := contexts.WithRequestRoute(request.Context(), route)
		next.ServeHTTP(w, request.WithContext(ctx))
	}
}
