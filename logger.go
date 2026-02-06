package log

import (
	"context"

	"go.uber.org/zap"
)

// Debug 记录Debug级别日志
// ctx: 上下文对象
// msg: 日志消息
// fields: 额外字段
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Debug(msg, logger.fields(ctx, fields...)...)
}

// Debugw 记录Debug级别日志（键值对形式）
// ctx: 上下文对象
// msg: 日志消息
// args: 键值对参数
func Debugw(ctx context.Context, msg string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Debugw(msg, args...)
}

// Info 记录Info级别日志
// ctx: 上下文对象
// msg: 日志消息
// fields: 额外字段
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(msg, logger.fields(ctx, fields...)...)
}

// Infow 记录Info级别日志（键值对形式）
// ctx: 上下文对象
// msg: 日志消息
// args: 键值对参数
func Infow(ctx context.Context, msg string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Infow(msg, args...)
}

// Warn 记录Warn级别日志
// ctx: 上下文对象
// msg: 日志消息
// fields: 额外字段
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Warn(msg, logger.fields(ctx, fields...)...)
}

// Warnw 记录Warn级别日志（键值对形式）
// ctx: 上下文对象
// msg: 日志消息
// args: 键值对参数
func Warnw(ctx context.Context, msg string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Warnw(msg, args...)
}

// Error 记录Error级别日志
// ctx: 上下文对象
// msg: 日志消息
// fields: 额外字段
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Error(msg, logger.fields(ctx, fields...)...)
}

// Errorw 记录Error级别日志（键值对形式）
// ctx: 上下文对象
// msg: 日志消息
// args: 键值对参数
func Errorw(ctx context.Context, msg string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Errorw(msg, args...)
}

// Fatal 记录Fatal级别日志并退出程序
// ctx: 上下文对象
// msg: 日志消息
// fields: 额外字段
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Fatal(msg, logger.fields(ctx, fields...)...)
}

// Panic 记录Panic级别日志并触发panic
// ctx: 上下文对象
// msg: 日志消息
// fields: 额外字段
func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Panic(msg, logger.fields(ctx, fields...)...)
}

// With 添加固定字段到日志对象
// fields: 固定字段
// return: 带有固定字段的日志对象
func With(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

// WithOptions 添加选项到日志对象
// fields: 日志选项
// return: 带有选项的日志对象
func WithOptions(fields ...zap.Option) *zap.Logger {
	return logger.WithOptions(fields...)
}

// Sync 刷新日志缓冲区
func Sync() {
	_ = logger.Sync()
}

// Debugf 记录Debug级别日志（格式化字符串）
// ctx: 上下文对象
// template: 格式化模板
// args: 格式化参数
func Debugf(ctx context.Context, template string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Debugf(template, args...)
}

// Infof 记录Info级别日志（格式化字符串）
// ctx: 上下文对象
// template: 格式化模板
// args: 格式化参数
func Infof(ctx context.Context, template string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Infof(template, args...)
}

// WithInfof 记录Info级别日志（带固定字段和格式化字符串）
// ctx: 上下文对象
// field: 固定字段
// template: 格式化模板
// args: 格式化参数
func WithInfof(ctx context.Context, field zap.Field, template string, args ...interface{}) {
	// 直接将 field 传递给 fields 方法，它会处理合并
	logger.With(logger.fields(ctx, field)...).Sugar().Infof(template, args...)
}

// Warnf 记录Warn级别日志（格式化字符串）
// ctx: 上下文对象
// template: 格式化模板
// args: 格式化参数
func Warnf(ctx context.Context, template string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Warnf(template, args...)
}

// Errorf 记录Error级别日志（格式化字符串）
// ctx: 上下文对象
// template: 格式化模板
// args: 格式化参数
func Errorf(ctx context.Context, template string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Errorf(template, args...)
}

// Fatalf 记录Fatal级别日志（格式化字符串）并退出程序
// ctx: 上下文对象
// template: 格式化模板
// args: 格式化参数
func Fatalf(ctx context.Context, template string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Fatalf(template, args...)
}

// Panicf 记录Panic级别日志（格式化字符串）并触发panic
// ctx: 上下文对象
// template: 格式化模板
// args: 格式化参数
func Panicf(ctx context.Context, template string, args ...interface{}) {
	logger.With(logger.fields(ctx)...).Sugar().Panicf(template, args...)
}
