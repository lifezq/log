// Copyright 2026 The Goutils Author. All Rights Reserved.
//
// -------------------------------------------------------------------

package log

import (
	"context"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options 日志配置结构体
type Options struct {
	CtxFields    []string      // 上下文中需要传递的字段键名列表
	Filename     string        // 日志文件名称
	MaxCount     uint          // 日志文件最大保存个数
	CallerEnable bool          // 是否开启调用栈信息
	LogLevel     zapcore.Level // 日志级别
	CloseConsole bool          // 是否关闭控制台输出
	Fields       []zap.Field   // 全局固定字段
}

// Logger 封装 zap.Logger 和上下文配置
type Logger struct {
	*zap.Logger
	ctxFields []string
}

// logger 全局日志实例
var logger = &Logger{
	Logger:    zap.NewNop(),
	ctxFields: nil,
}

// fields 从上下文中提取配置的字段
func (l *Logger) fields(ctx context.Context, fields ...zap.Field) []zap.Field {
	if len(l.ctxFields) == 0 {
		return fields
	}

	result := make([]zap.Field, 0, len(l.ctxFields)+len(fields))
	for _, field := range l.ctxFields {
		val := ctx.Value(field)
		if val != nil {
			if strVal, ok := val.(string); ok {
				result = append(result, zap.String(field, strVal))
			}
		}
	}

	if len(fields) > 0 {
		result = append(result, fields...)
	}

	return result
}

// Init 初始化日志对象
// opt: 日志配置选项
func Init(opt Options) {
	var zLogger *zap.Logger

	// 配置日志编码器
	cfg := zapcore.EncoderConfig{
		MessageKey:   "message",
		TimeKey:      "timestamp",
		LevelKey:     "level",
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format(time.DateTime))
		},
	}

	// 根据配置创建日志实例
	if opt.CloseConsole {
		zLogger = zap.New(
			zapcore.NewTee(
				syncFile(cfg, opt),
			),
		)
	} else {
		zLogger = zap.New(
			zapcore.NewTee(
				syncConsole(cfg, opt),
				syncFile(cfg, opt),
			),
		)
	}

	// 开启调用栈信息
	if opt.CallerEnable {
		zLogger = zLogger.WithOptions(
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)
	}

	// 添加全局固定字段
	if len(opt.Fields) > 0 {
		zLogger = zLogger.With(opt.Fields...)
	}

	// 更新全局实例
	logger = &Logger{
		Logger:    zLogger,
		ctxFields: opt.CtxFields,
	}
}

// syncFile 创建文件输出的日志核心
// cfg: 编码器配置
// opt: 日志配置选项
// return: 日志核心对象
func syncFile(cfg zapcore.EncoderConfig, opt Options) zapcore.Core {
	// 创建日志文件轮转器
	writer, err := rotatelogs.New(
		opt.Filename+".%Y%m%d.log",
		rotatelogs.WithLinkName(opt.Filename),
		rotatelogs.WithRotationCount(opt.MaxCount),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}

	// 创建并返回日志核心
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(writer),
		opt.LogLevel,
	)
}

// syncConsole 创建控制台输出的日志核心
// cfg: 编码器配置
// opt: 日志配置选项
// return: 日志核心对象
func syncConsole(cfg zapcore.EncoderConfig, opt Options) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		opt.LogLevel,
	)
}
