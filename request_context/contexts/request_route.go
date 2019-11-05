package contexts

import (
	"context"
)

type requestRouteContext int

const requestRouteKey requestRouteContext = 0

func GetRequestRoute(ctx context.Context) string {
	if requestRoute, ok := ctx.Value(requestRouteKey).(string); ok {
		return requestRoute
	}
	return ""
}

func WithRequestRoute(ctx context.Context, requestRoute string) context.Context {
	return context.WithValue(ctx, requestRouteKey, requestRoute)
}
