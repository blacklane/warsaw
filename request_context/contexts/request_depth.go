package contexts

import (
	"context"
)

type requestDepthContext int

const requestDepthKey requestDepthContext = 0

func GetRequestDepth(ctx context.Context) int {
	if requestDepth, ok := ctx.Value(requestDepthKey).(int); ok {
		return requestDepth
	}
	return 0
}

func WithRequestDepth(ctx context.Context, requestDepth int) context.Context {
	return context.WithValue(ctx, requestDepthKey, requestDepth)
}
