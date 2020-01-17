package request_context

import (
	"context"
	"net/http"
	"strconv"

	"github.com/blacklane/warsaw/request_context/constants"
	"github.com/blacklane/warsaw/request_context/contexts"
)

// TrackerMiddleware will record the `RequestContext` instance in the context of the request + pass the RequestId
// to the response Headers accordingly.
func TrackerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		requestContext := ExtractRequestContext(request)
		ctx := buildContextFromRequestContext(request.Context(), requestContext)

		w.Header().Set(constants.RequestIDHeader, requestContext.RequestId)
		next.ServeHTTP(w, request.WithContext(ctx))
	}
}

// SetTrackerHeaders is useful if you want to pass headers to a downstream net/http.Request.
func SetTrackerHeaders(ctx context.Context, header *http.Header) {
	header.Set(constants.RequestIDHeader, contexts.GetRequestID(ctx))
	header.Set(constants.RequestDepthHeader, strconv.Itoa(contexts.GetRequestDepth(ctx)))
	header.Set(constants.TreePathHeader, contexts.GetTreePath(ctx))
}

func buildContextFromRequestContext(ctx context.Context, requestContext RequestContext) context.Context {
	ctx = contexts.WithRequestID(ctx, requestContext.RequestId)
	ctx = contexts.WithRequestDepth(ctx, requestContext.RequestDepth)
	ctx = contexts.WithTreePath(ctx, requestContext.TreePath)
	return ctx
}
