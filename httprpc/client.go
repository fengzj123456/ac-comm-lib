package httprpc

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/httprpc/codes"
)

var HTTPClient = http.Client{
	Timeout: 20 * time.Second,
}

type Client struct {
	url   string
	codec Codec
}

func NewClient(url string, codec Codec) *Client {
	if codec == nil {
		codec = DefaultCodec
	}
	return &Client{url: url, codec: codec}
}

func (c *Client) Call(ctx context.Context, path string, args, reply interface{}) (err error) {
	var buf bytes.Buffer
	if args != nil {
		if err = c.codec.Encode(&buf, args); err != nil {
			return err
		}
	}

	url := c.url + normalizePath(path)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return err
	}
	c.setRequestHeader(req, ctx)

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var er errReply
		if err = c.codec.Decode(resp.Body, &er); err != nil {
			return err
		}
		return clientError(codes.Code(er.Code), er.Cause, er.Stack)
	}
	if reply != nil {
		if err = c.codec.Decode(resp.Body, reply); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) setRequestHeader(r *http.Request, ctx context.Context) {
	setHeaderContentType(r.Header, c.codec.ContentType())
	if rctx, ok := ctx.(*Context); ok {
		setHeaderTraceID(r.Header, rctx.TraceID)
		for key, values := range rctx.RequestHeader {
			for _, value := range values {
				r.Header.Add(key, value)
			}
		}
	} else {
		setHeaderTraceID(r.Header, "")
	}
}
