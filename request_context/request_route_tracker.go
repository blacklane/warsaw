package request_context

import (
	"net/http"

	"github.com/blacklane/warsaw/request_context/contexts"
)

// RequestRouteTracker stores the routeName in the logger context for future lookups
func RequestRouteTracker(route string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := contexts.WithRequestRoute(request.Context(), route)
		next.ServeHTTP(w, request.WithContext(ctx))
	}
}
