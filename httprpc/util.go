package httprpc

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/ironzhang/pearls/uuid"
)

const (
	xTraceID = "X-Trace-Id"
)

func normalizePath(path string) string {
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		path = "/"
	} else if path[0] != '/' {
		path = "/" + path
	}
	return path
}

func splitPath(path string) (string, string) {
	path = normalizePath(path)
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return "/", path
	} else if i == 0 {
		return path[:1], path[1:]
	}
	return path[:i], path[i+1:]
}

func runtimeStack() []byte {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	return buf
}

func getHeaderContentType(h http.Header) string {
	return h.Get("Content-Type")
}

func setHeaderContentType(h http.Header, contentType string) {
	h.Set("Content-Type", contentType)
}

func getHeaderTraceID(h http.Header) string {
	if traceID := h.Get(xTraceID); traceID != "" {
		return traceID
	}
	return uuid.New().String()
}

func setHeaderTraceID(h http.Header, traceID string) {
	if traceID == "" {
		traceID = uuid.New().String()
	}
	h.Set(xTraceID, traceID)
}

func setCors(h http.Header, origin string) {
	h.Set("Access-Control-Allow-Origin", origin)
	h.Set("Access-Control-Allow-Methods", "*")
	h.Set("Access-Control-Allow-Headers", "*")
}
