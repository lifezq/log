package log

import (
	"context"
	"io"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestInitLog 测试日志初始化
func TestInitLog(t *testing.T) {
	// 测试初始化日志
	InitLog(Options{
		Filename:     "test",
		MaxCount:     1,
		CallerEnable: true,
		LogLevel:     zapcore.DebugLevel,
		CloseConsole: false,
	}, zap.String("test", "value"))

	// 测试Sync函数
	Sync()
}

// TestBasicLog 测试基本日志功能
func TestBasicLog(t *testing.T) {
	ctx := context.Background()

	// 测试各种级别的日志
	Debug(ctx, "Debug message")
	Info(ctx, "Info message")
	Warn(ctx, "Warn message")
	Error(ctx, "Error message")

	// 测试格式化日志
	Debugf(ctx, "Debug formatted message: %s", "test")
	Infof(ctx, "Info formatted message: %s", "test")
	Warnf(ctx, "Warn formatted message: %s", "test")
	Errorf(ctx, "Error formatted message: %s", "test")

	// 测试键值对日志
	Debugw(ctx, "Debug key-value message", "key", "value")
	Infow(ctx, "Info key-value message", "key", "value")
	Warnw(ctx, "Warn key-value message", "key", "value")
	Errorw(ctx, "Error key-value message", "key", "value")
}

// TestContextLog 测试上下文日志功能
func TestContextLog(t *testing.T) {
	// 重新初始化日志，配置上下文字段
	InitLog(Options{
		Filename:     "test",
		MaxCount:     1,
		CallerEnable: true,
		LogLevel:     zapcore.DebugLevel,
		CloseConsole: false,
		CtxFields:    []string{"request_id", "user_id"},
	})

	// 创建带有值的上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "123456")
	ctx = context.WithValue(ctx, "user_id", "789")

	// 测试上下文日志
	Info(ctx, "Context log message")
	Error(ctx, "Context error message")
}

// TestWithFunctions 测试With相关函数
func TestWithFunctions(t *testing.T) {
	// 测试With函数
	logger := With(zap.String("test", "value"))
	if logger == nil {
		t.Errorf("With function returned nil")
	}

	// 测试WithOptions函数
	loggerWithOptions := WithOptions(zap.AddCaller())
	if loggerWithOptions == nil {
		t.Errorf("WithOptions function returned nil")
	}
}

// TestSafeWriter 测试SafeWriter功能
func TestSafeWriter(t *testing.T) {
	ctx := context.Background()

	// 获取SafeWriter
	writer := SafeWriter(ctx)
	if writer == nil {
		t.Errorf("SafeWriter returned nil")
	}

	// 写入测试数据
	testData := "Test SafeWriter data"
	_, err := writer.Write([]byte(testData))
	if err != nil {
		t.Errorf("Error writing to SafeWriter: %v", err)
	}

	// 关闭writer
	err = writer.Close()
	if err != nil {
		t.Errorf("Error closing SafeWriter: %v", err)
	}

	// 等待写入完成
	time.Sleep(100 * time.Millisecond)
}

// TestGormLogger 测试GORM日志记录器
func TestGormLogger(t *testing.T) {
	ctx := context.Background()

	// 创建GormLogger
	gormLogger := GormLogger{
		LogLevel:                  1, // Info级别
		IgnoreRecordNotFoundError: true,
		SlowThreshold:             100 * time.Millisecond,
	}

	// 测试LogMode方法
	loggerInterface := gormLogger.LogMode(2) // Warn级别
	if loggerInterface == nil {
		t.Errorf("LogMode returned nil")
	}

	// 测试Info方法
	gormLogger.Info(ctx, "Info message", "param")

	// 测试Warn方法
	gormLogger.Warn(ctx, "Warn message", "param")

	// 测试Error方法
	gormLogger.Error(ctx, "Error message", "param")

	// 测试Trace方法
	begin := time.Now()
	gormLogger.Trace(ctx, begin, func() (string, int64) {
		return "SELECT * FROM users", 1
	}, nil)
}

// TestEsLogger 测试Elasticsearch日志记录器
func TestEsLogger(t *testing.T) {
	// 创建EsLogger
	esLogger := EsLogger{
		RequestEnabled:  true,
		ResponseEnabled: true,
	}

	// 测试RequestBodyEnabled方法
	if !esLogger.RequestBodyEnabled() {
		t.Errorf("RequestBodyEnabled should return true")
	}

	// 测试ResponseBodyEnabled方法
	if !esLogger.ResponseBodyEnabled() {
		t.Errorf("ResponseBodyEnabled should return true")
	}
}

// TestScanLinesOrGiveLong 测试scanLinesOrGiveLong函数
func TestScanLinesOrGiveLong(t *testing.T) {
	// 测试正常行
	data := []byte("test\n")
	advance, token, err := scanLinesOrGiveLong(data, false)
	if advance == 0 || len(token) == 0 || err != nil {
		t.Errorf("scanLinesOrGiveLong failed for normal line")
	}

	// 测试长数据
	longData := make([]byte, maxTokenLength+1)
	for i := range longData {
		longData[i] = 'a'
	}
	advance, token, err = scanLinesOrGiveLong(longData, false)
	if advance != maxTokenLength || len(token) != maxTokenLength || err != nil {
		t.Errorf("scanLinesOrGiveLong failed for long data")
	}
}

// TestWriterFinalizer 测试writerFinalizer函数
func TestWriterFinalizer(t *testing.T) {
	// 创建管道
	reader, writer := io.Pipe()

	// 测试writerFinalizer
	writerFinalizer(writer)

	// 关闭reader
	_ = reader.Close()
}
