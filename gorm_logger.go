// Copyright 2026 The Goutils Author. All Rights Reserved.
//
// -------------------------------------------------------------------

package log

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

// GormLogger GORM日志记录器
type GormLogger struct {
	LogLevel                  gl.LogLevel   // 日志级别
	IgnoreRecordNotFoundError bool          // 是否忽略记录未找到错误
	SlowThreshold             time.Duration // 慢SQL阈值
}

// LogMode 设置日志模式
// level: 日志级别
// return: 日志接口
func (l GormLogger) LogMode(level gl.LogLevel) gl.Interface {
	return &GormLogger{
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

// Info 记录Info级别日志
// ctx: 上下文对象
// str: 日志消息模板
// args: 日志消息参数
func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gl.Info {
		return
	}
	Infof(ctx, str, args...)
}

// Warn 记录Warn级别日志
// ctx: 上下文对象
// str: 日志消息模板
// args: 日志消息参数
func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gl.Warn {
		return
	}
	Warnf(ctx, str, args...)
}

// Error 记录Error级别日志
// ctx: 上下文对象
// str: 日志消息模板
// args: 日志消息参数
func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gl.Error {
		return
	}
	Errorf(ctx, str, args...)
}

// Trace 记录SQL执行轨迹
// ctx: 上下文对象
// begin: 开始时间
// fc: 获取SQL和行数的函数
// err: 错误信息
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gl.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		// 记录错误日志
		sql, rows := fc()
		Errorf(ctx, "err=%s elapsed=%s rows=%d sql=%s", err.Error(), elapsed.String(), rows, sql)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gl.Warn:
		// 记录慢SQL日志
		sql, rows := fc()
		var e string
		if err != nil {
			e = err.Error()
		}
		Warnf(ctx, "err=%s elapsed=%s rows=%d sql=%s", e, elapsed.String(), rows, sql)
	case l.LogLevel >= gl.Info:
		// 记录普通SQL日志
		sql, rows := fc()
		var e string
		if err != nil {
			e = err.Error()
		}
		Infof(ctx, "err=%s elapsed=%s rows=%d sql=%s", e, elapsed.String(), rows, sql)
	}
}
