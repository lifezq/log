[![Build status](https://img.shields.io/appveyor/build/lifezq/log.svg)](https://ci.appveyor.com/project/lifezq/log)
[![Coverage Status](https://img.shields.io/coveralls/lifezq/log.svg?style=flat-square)](https://coveralls.io/github/lifezq/log?branch=master)
[![License](http://img.shields.io/badge/license-apache-blue.svg?style=flat-square)](https://raw.githubusercontent.com/lifezq/log/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/lifezq/log)

# log

一个功能丰富、易于使用的Go语言日志库，基于zap封装，提供了上下文日志、文件轮转、多级别日志、以及与GORM和Elasticsearch的集成支持。

## 功能特性

- ✅ 基于zap的高性能日志实现
- ✅ 支持上下文日志，可从上下文中提取字段
- ✅ 支持日志文件轮转
- ✅ 支持多级别日志（Debug、Info、Warn、Error、Fatal、Panic）
- ✅ 支持控制台和文件双输出
- ✅ 支持与GORM集成
- ✅ 支持与Elasticsearch集成
- ✅ 提供io.Writer接口，方便与其他库集成
- ✅ 支持结构化日志（JSON格式）

## 安装

```bash
go get github.com/lifezq/log
```

## 快速开始

### 基本使用

```go
package main

import (
    "context"
    "github.com/lifezq/log"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    // 初始化日志
    log.Init(log.Options{
        Filename:     "app",           // 日志文件名称
        MaxCount:     7,               // 最大保存7个日志文件
        CallerEnable: true,            // 开启调用栈信息
        LogLevel:     zapcore.InfoLevel, // 日志级别
        CloseConsole: false,           // 不关闭控制台输出
        Fields: []zap.Field{           // 全局固定字段
            zap.String("service", "example"),
        },
    })

    // 创建上下文
    ctx := context.Background()
    ctx = context.WithValue(ctx, "request_id", "123456")

    // 记录日志
    log.Infof(ctx, "Hello, world!")
    log.Errorf(ctx, "Error occurred: %v", "something went wrong")

    // 刷新日志缓冲区
    log.Sync()
}
```

## 配置选项

| 配置项 | 类型 | 描述 | 默认值 |
|-------|------|------|-------|
| Filename | string | 日志文件名称 | - |
| MaxCount | uint | 日志文件最大保存个数 | - |
| CallerEnable | bool | 是否开启调用栈信息 | false |
| LogLevel | zapcore.Level | 日志级别 | zapcore.InfoLevel |
| CloseConsole | bool | 是否关闭控制台输出 | false |
| CtxFields | []string | 上下文中需要传递的字段键名列表 | [] |
| Fields | []zap.Field | 全局固定字段 | [] |

## 日志级别

- `zapcore.DebugLevel` - 调试级别
- `zapcore.InfoLevel` - 信息级别
- `zapcore.WarnLevel` - 警告级别
- `zapcore.ErrorLevel` - 错误级别
- `zapcore.DPanicLevel` - 开发模式panic
- `zapcore.PanicLevel` - panic级别
- `zapcore.FatalLevel` - 致命错误级别

## 使用示例

### 1. 基本日志

```go
// 记录Info级别日志
log.Info(ctx, "This is an info message")

// 记录Error级别日志
log.Error(ctx, "This is an error message", zap.String("detail", "error details"))

// 使用格式化字符串
log.Infof(ctx, "Hello %s!", "world")

// 使用键值对形式
log.Infow(ctx, "User login", "user_id", 123, "ip", "192.168.1.1")
```

### 2. 上下文日志

```go
// 在上下文中设置值
ctx := context.Background()
ctx = context.WithValue(ctx, "request_id", "123456")
ctx = context.WithValue(ctx, "user_id", 789)

// 初始化时配置要从上下文中提取的字段
    log.Init(log.Options{
        // ... 其他配置
        CtxFields: []string{"request_id", "user_id"}, // 从上下文中提取这些字段
        Fields: []zap.Field{
            zap.String("service", "api"),
        },
    })

// 记录日志，会自动包含上下文中的字段
log.Info(ctx, "API request received")
// 输出会包含 request_id=123456 和 user_id=789
```

### 3. 与GORM集成

```go
import (
    "github.com/lifezq/log"
    "gorm.io/gorm"
)

func setupGORM() *gorm.DB {
    // 初始化GORM日志记录器
    gormLogger := log.GormLogger{
        LogLevel:                  gl.Info,           // 日志级别
        IgnoreRecordNotFoundError: true,              // 忽略记录未找到错误
        SlowThreshold:             200 * time.Millisecond, // 慢SQL阈值
    }

    // 连接数据库时使用自定义日志记录器
    db, err := gorm.Open(mysql.Open("dsn"), &gorm.Config{
        Logger: gormLogger.LogMode(gl.Info),
    })
    
    return db
}
```

### 4. 与Elasticsearch集成

```go
import (
    "github.com/lifezq/log"
    "github.com/elastic/go-elasticsearch/v7"
)

func setupES() (*elasticsearch.Client, error) {
    // 初始化Elasticsearch客户端
    cfg := elasticsearch.Config{
        Addresses: []string{"http://localhost:9200"},
        Logger: log.EsLogger{
            RequestEnabled:  true,  // 开启请求日志
            ResponseEnabled: false, // 关闭响应日志
        },
    }

    client, err := elasticsearch.NewClient(cfg)
    return client, err
}
```

### 5. 使用io.Writer接口

```go
import (
    "context"
    "fmt"
    "github.com/lifezq/log"
)

func main() {
    ctx := context.Background()
    
    // 获取一个io.Writer，写入的内容会被记录为Info级别日志
    writer := log.SafeWriter(ctx)
    
    // 使用fmt.Fprintf写入内容
    fmt.Fprintf(writer, "This will be logged as info level")
    fmt.Fprintf(writer, "\nAnother line to log")
}
```

## API文档

### 核心函数

#### `Init(opt Options)`
初始化日志对象
- `opt`: 日志配置选项，包含全局固定字段

#### 日志记录函数

- `Debug(ctx context.Context, msg string, fields ...zap.Field)` - 记录Debug级别日志
- `Info(ctx context.Context, msg string, fields ...zap.Field)` - 记录Info级别日志
- `Warn(ctx context.Context, msg string, fields ...zap.Field)` - 记录Warn级别日志
- `Error(ctx context.Context, msg string, fields ...zap.Field)` - 记录Error级别日志
- `Fatal(ctx context.Context, msg string, fields ...zap.Field)` - 记录Fatal级别日志并退出程序
- `Panic(ctx context.Context, msg string, fields ...zap.Field)` - 记录Panic级别日志并触发panic

#### 格式化日志函数

- `Debugf(ctx context.Context, template string, args ...interface{})` - 记录Debug级别格式化日志
- `Infof(ctx context.Context, template string, args ...interface{})` - 记录Info级别格式化日志
- `Warnf(ctx context.Context, template string, args ...interface{})` - 记录Warn级别格式化日志
- `Errorf(ctx context.Context, template string, args ...interface{})` - 记录Error级别格式化日志
- `Fatalf(ctx context.Context, template string, args ...interface{})` - 记录Fatal级别格式化日志并退出程序
- `Panicf(ctx context.Context, template string, args ...interface{})` - 记录Panic级别格式化日志并触发panic

#### 键值对日志函数

- `Debugw(ctx context.Context, msg string, args ...interface{})` - 记录Debug级别键值对日志
- `Infow(ctx context.Context, msg string, args ...interface{})` - 记录Info级别键值对日志
- `Warnw(ctx context.Context, msg string, args ...interface{})` - 记录Warn级别键值对日志
- `Errorw(ctx context.Context, msg string, args ...interface{})` - 记录Error级别键值对日志

#### 其他函数

- `Sync()` - 刷新日志缓冲区
- `With(fields ...zap.Field)` - 添加固定字段到日志对象
- `WithOptions(fields ...zap.Option)` - 添加选项到日志对象
- `SafeWriter(ctx context.Context)` - 获取一个io.Writer，写入的内容会被记录为日志

## 配置示例

### 生产环境配置

```go
log.Init(log.Options{
    Filename:     "/var/log/app/app",  // 日志文件路径
    MaxCount:     30,                  // 保存30天日志
    CallerEnable: false,               // 生产环境关闭调用栈信息
    LogLevel:     zapcore.WarnLevel,   // 只记录Warn及以上级别
    CloseConsole: true,                // 关闭控制台输出
    Fields: []zap.Field{
        zap.String("env", "production"),
    },
})
```

### 开发环境配置

```go
log.Init(log.Options{
    Filename:     "app",              // 日志文件名称
    MaxCount:     7,                   // 保存7天日志
    CallerEnable: true,                // 开启调用栈信息
    LogLevel:     zapcore.DebugLevel,  // 记录所有级别日志
    CloseConsole: false,               // 开启控制台输出
    Fields: []zap.Field{
        zap.String("env", "development"),
    },
})
```

## 许可证

Apache License
