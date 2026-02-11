// Copyright 2026 The Goutils Author. All Rights Reserved.
//
// -------------------------------------------------------------------

package log

import (
	"bufio"
	"context"
	"io"
	"runtime"
	"strings"
)

// SafeWriter 创建一个安全的io.Writer，将写入的内容作为日志记录
// ctx: 上下文对象
// return: 管道写入器
func SafeWriter(ctx context.Context) *io.PipeWriter {
	reader, writer := io.Pipe()
	go scan(ctx, reader)
	runtime.SetFinalizer(writer, writerFinalizer)
	return writer
}

// scan 扫描管道读取器中的内容并记录为日志
// ctx: 上下文对象
// reader: 管道读取器
func scan(ctx context.Context, reader *io.PipeReader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesOrGiveLong)
	for scanner.Scan() {
		// 跳过空行
		if strings.Replace(scanner.Text(), " ", "", -1) == "" {
			continue
		}

		Infof(ctx, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		Errorf(ctx, "Error while reading from Writer: %s", err)
	}
	_ = reader.Close()
}

const maxTokenLength = bufio.MaxScanTokenSize / 2

// scanLinesOrGiveLong 自定义扫描函数，处理长行
// data: 数据
// atEOF: 是否到达文件末尾
// return: 前进长度, 令牌, 错误
func scanLinesOrGiveLong(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	if advance > 0 || token != nil || err != nil {
		return
	}
	// 如果数据长度超过最大令牌长度，返回部分数据
	if len(data) < maxTokenLength {
		return
	}
	return maxTokenLength, data[0:maxTokenLength], nil
}

// writerFinalizer 管道写入器的终结器
// writer: 管道写入器
func writerFinalizer(writer *io.PipeWriter) {
	_ = writer.Close()
}
