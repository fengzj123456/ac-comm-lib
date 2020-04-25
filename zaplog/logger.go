package zaplog

import (
	"context"

	"go.uber.org/zap"
)

type TraceIDContext interface {
	GetTraceID() string
}

type values struct {
	traceID  string
	errorStr string
	function string
}

type Logger struct {
	base   *zap.Logger
	values values
	fields []zap.Field
}

func NewLogger(base *zap.Logger) *Logger {
	return &Logger{
		base: base.WithOptions(zap.AddCallerSkip(1)),
	}
}

func (l *Logger) clone(n int) *Logger {
	c := &Logger{
		base:   l.base,
		values: l.values,
	}
	c.fields = make([]zap.Field, len(l.fields), len(l.fields)+n)
	copy(c.fields, l.fields)
	return c
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	c := l.clone(len(fields))
	c.fields = append(c.fields, fields...)
	return c
}

func (l *Logger) WithArgs(args ...interface{}) *Logger {
	return l.With(l.sweetenFields(args)...)
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	log := l
	if rc, ok := ctx.(TraceIDContext); ok {
		if id := rc.GetTraceID(); id != "" {
			log = log.clone(0)
			log.values.traceID = id
		}
	}
	return log
}

func (l *Logger) WithFunction(name string) *Logger {
	log := l
	if name != "" {
		log = log.clone(0)
		log.values.function = name
	}
	return log
}

func (l *Logger) Base() *zap.Logger {
	log := l.base
	if l.values.traceID != "" {
		log = log.With(zap.String("TraceID", l.values.traceID))
	}
	if l.values.errorStr != "" {
		log = log.With(zap.String("Error", l.values.errorStr))
	}
	if l.values.function != "" {
		log = log.With(zap.String("Function", l.values.function))
	}
	if len(l.fields) > 0 {
		log = log.With(l.fields...)
	}
	return log
}

func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.Base().Sugar()
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Base().Debug(msg, fields...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Debugw(msg, keysAndValues...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.Sugar().Debugf(template, args...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Base().Info(msg, fields...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.Sugar().Infof(template, args...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Base().Warn(msg, fields...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Warnw(msg, keysAndValues...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.Sugar().Warnf(template, args...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Base().Error(msg, fields...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Errorw(msg, keysAndValues...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.Sugar().Errorf(template, args...)
}

func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.Base().DPanic(msg, fields...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.Sugar().DPanicw(msg, keysAndValues...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.Sugar().DPanicf(template, args...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.Base().Panic(msg, fields...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Panicw(msg, keysAndValues...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.Sugar().Panicf(template, args...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Base().Fatal(msg, fields...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.Sugar().Fatalf(template, args...)
}
