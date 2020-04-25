package zaplog

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

type TContext struct {
	context.Context
	id string
}

func (c *TContext) GetTraceID() string {
	return c.id
}

func NewTestLogger(t *testing.T) *Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("new development: %v", err)
	}
	return NewLogger(l)
}

func TestLogger(t *testing.T) {
	logger := NewTestLogger(t).WithFunction("TestLogger")
	logger.Debug("debug", zap.String("A", "a"), zap.Int64("B", 1))
	logger.Debugw("debugw", "A", "a", "B", 1)
	logger.Debugf("debugf A=%s, B=%d", "a", 1)
	logger.WithContext(&TContext{id: "1"}).Debug("debug", zap.String("A", "a"), zap.Int64("B", 1))
	logger.WithContext(&TContext{id: "2"}).Debugw("debugw", "A", "a", "B", 1)
	logger.WithContext(&TContext{id: "3"}).Debugf("debugf A=%s, B=%d", "a", 1)
	logger.WithContext(NewContext()).Debugf("context")
}
