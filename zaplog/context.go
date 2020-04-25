package zaplog

import (
	"context"

	"github.com/ironzhang/pearls/uuid"
)

type Context struct {
	context.Context
	TraceID string
}

func (c *Context) GetTraceID() string {
	return c.TraceID
}

func NewContext() *Context {
	return &Context{
		Context: context.Background(),
		TraceID: uuid.New().String(),
	}
}
