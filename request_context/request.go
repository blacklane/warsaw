package request_context

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/blacklane/warsaw/constants"
)

// RequestContext provides details on the HTTP request tracked by warsaw/Logger
type RequestContext struct {
	RequestID, TreePath string
	RequestDepth        int
}

// ExtractRequestContext creates the RequestContext from http.Request instance
func ExtractRequestContext(r *http.Request) RequestContext {
	return RequestContext{
		extractRequestID(r),
		extractTreePath(r),
		extractRequestDepth(r),
	}
}

func extractRequestID(r *http.Request) string {
	requestID := r.Header.Get(constants.HeaderRequestID)
	if len(requestID) == 0 {
		requestID = uuid.New().String()
	}
	return requestID
}

func extractRequestDepth(r *http.Request) int {
	depth, err := strconv.Atoi(r.Header.Get(constants.HeaderRequestDepth))
	if err != nil {
		depth = 0
	} else {
		depth++
	}
	return depth
}

// TreePathSuffix defines what should be added to Request tracking tree_path as definition of the current path.
var TreePathSuffix = "T"

func extractTreePath(r *http.Request) string {
	treePath := r.Header.Get(constants.HeaderTreePath) + TreePathSuffix

	return treePath
}
