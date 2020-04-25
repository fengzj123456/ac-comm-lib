package httprpc

import (
	"context"
	"net/http"
)

type Context struct {
	context.Context
	TraceID        string
	Request        *http.Request
	Response       http.ResponseWriter
	RequestHeader  http.Header // client使用, 该Header中的值将在请求发送时添加到Request中.
	ResponseHeader http.Header // server使用, 该Header中的值将在响应返回时添加到Response中.
}

func (c *Context) GetTraceID() string {
	return c.TraceID
}
