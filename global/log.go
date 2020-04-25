package global

import (
	"context"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/zaplog"
)

func Logger(ctx context.Context) *zaplog.Logger {
	return zaplog.Std.WithContext(ctx)
}

func SugaredLogger(ctx context.Context) *zaplog.Logger {
	return zaplog.Std.WithContext(ctx)
}
