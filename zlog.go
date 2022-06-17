// Package zlog is zap in context wrapper
package zlog

import (
	"context"

	"go.uber.org/zap"
)

type contextKey int

const (
	ctxLoggerKey contextKey = iota + 1
	ctxFieldsKey
)

var (
	// CtxWithLogger is set logger in context.Context
	CtxWithLogger func(ctx context.Context, logger *zap.Logger) context.Context = ctxWithLoggerFun
	// Logger return the zap.Logger from context.
	// if the context don't have logger then return the zap.L()
	Logger func(ctx context.Context) *zap.Logger = loggerFun
)

// SetGlobalLogger is replace global zap.Logger (default logger is Nop)
// return the function revert before
func SetGlobalLogger(logger *zap.Logger) func() {
	return zap.ReplaceGlobals(logger)
}

func ctxWithLoggerFun(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

func loggerFun(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxLoggerKey).(*zap.Logger); ok {
		return logger.With(Fields(ctx)...)
	}
	return zap.L().With(Fields(ctx)...)
}

// With set zap.Field to context
func With(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, ctxFieldsKey, MergeFields(Fields(ctx), fields))
}

// Fields return zap.Fields in context
func Fields(ctx context.Context) []zap.Field {
	if fields, ok := ctx.Value(ctxFieldsKey).([]zap.Field); ok {
		return fields
	}
	return []zap.Field{}
}

// MergeFields はbaseをoverrideで上書きしてmergeする
func MergeFields(base, override []zap.Field) []zap.Field {
	overrideKeys := make(map[string]struct{}, len(override))
	for _, f := range override {
		overrideKeys[f.Key] = struct{}{}
	}
	after := make([]zap.Field, 0, len(base)+len(override))
	for _, f := range base {
		if _, exists := overrideKeys[f.Key]; !exists {
			after = append(after, f)
		}
	}
	return append(after, override...)
}

// Sugar is alias zap.Logger.Sugar()
func Sugar(ctx context.Context) *zap.SugaredLogger {
	return Logger(ctx).Sugar()
}

// Debugf is wrapper SugaredLogger.Debugf
func Debugf(ctx context.Context, template string, args ...interface{}) {
	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Sugar().Debugf(template, args...)
}

// Infof is wrapper SugaredLogger.Infof
func Infof(ctx context.Context, template string, args ...interface{}) {
	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Sugar().Infof(template, args...)
}

// Warnf is wrapper SugaredLogger.Warnf
func Warnf(ctx context.Context, template string, args ...interface{}) {
	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Sugar().Warnf(template, args...)
}

// Errorf is wrapper SugaredLogger.Errorf
func Errorf(ctx context.Context, template string, args ...interface{}) {
	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Sugar().Errorf(template, args...)
}

// Panicf is wrapper SugaredLogger.Panicf
func Panicf(ctx context.Context, template string, args ...interface{}) {
	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Sugar().Panicf(template, args...)
}
