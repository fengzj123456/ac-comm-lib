package zaplog_test

import (
	"context"
	"errors"
	"testing"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/zaplog"
)

type TContext struct {
	context.Context
	id string
}

func (c *TContext) GetTraceID() string {
	return c.id
}

func open(file string) error {
	return errors.New("failed to open")
}

func openFileA(file string) {
	if err := open(file); err != nil {
		zaplog.Std.Infow("open", "file", file, "error", err)
	}
}

func openFileB(ctx context.Context, file string) {
	if err := open(file); err != nil {
		zaplog.Std.WithContext(ctx).Infow("open", "file", file, "error", err)
	}
}

func openFileC(ctx context.Context, file string) {
	log := zaplog.Std.WithContext(ctx).WithFunction("openFileC").WithArgs("file", file)
	if err := open(file); err != nil {
		log.Infow("open", "error", err)
	}
}

func TestLog(t *testing.T) {
	openFileA("a")
	openFileB(context.Background(), "b")
	openFileC(&TContext{id: "1"}, "c")
}
