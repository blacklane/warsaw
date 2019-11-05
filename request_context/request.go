package request_context

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/blacklane/warsaw/request_context/constants"
)

type RequestContext struct {
	RequestId, TreePath string
	RequestDepth        int
}

func ExtractRequestContext(r *http.Request) RequestContext {
	return RequestContext{
		extractRequestID(r),
		extractTreePath(r),
		extractRequestDepth(r),
	}
}

func extractRequestID(r *http.Request) string {
	requestID := r.Header.Get(constants.RequestIDHeader)
	if len(requestID) == 0 {
		requestID = uuid.New().String()
	}
	return requestID
}

func extractRequestDepth(r *http.Request) int {
	depth, err := strconv.Atoi(r.Header.Get(constants.RequestDepthHeader))
	if err != nil {
		depth = 0
	} else {
		depth += 1
	}
	return depth
}

var TreePathSuffix = "T"

func extractTreePath(r *http.Request) string {
	treePath := r.Header.Get(constants.TreePathHeader) + TreePathSuffix

	return treePath
}
