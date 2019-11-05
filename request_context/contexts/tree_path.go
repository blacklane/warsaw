package contexts

import (
	"context"
)

type treePathContext int

const treePathKey treePathContext = 0

func GetTreePath(ctx context.Context) string {
	if treePath, ok := ctx.Value(treePathKey).(string); ok {
		return treePath
	}
	return ""
}

func WithTreePath(ctx context.Context, treePath string) context.Context {
	return context.WithValue(ctx, treePathKey, treePath)
}
