// Copyright 2026 The Goutils Author. All Rights Reserved.
//
// -------------------------------------------------------------------

package log

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"time"
)

// EsLogger Elasticsearch日志记录器
type EsLogger struct {
	RequestEnabled  bool `desc:"是否开启请求日志"` // 是否开启请求日志
	ResponseEnabled bool `desc:"是否开启响应日志"` // 是否开启响应日志
}

// LogRoundTrip 记录HTTP请求和响应
// request: HTTP请求
// response: HTTP响应
// err: 错误信息
// t: 请求时间
// ti: 请求耗时
// return: 错误信息
func (logger EsLogger) LogRoundTrip(request *http.Request, response *http.Response, err error, t time.Time, ti time.Duration) error {
	// 记录请求体
	if logger.RequestBodyEnabled() && request != nil && request.Body != nil && request.Body != http.NoBody {
		var buf bytes.Buffer
		if request.GetBody != nil {
			b, _ := request.GetBody()
			_, err = buf.ReadFrom(b)
		} else {
			_, err = buf.ReadFrom(request.Body)
		}
		if err != nil {
			Errorf(context.TODO(), "ES日志error: %+v", err)
			return nil
		}
		// 逐行扫描并记录请求体
		scanner := bufio.NewScanner(&buf)
		for scanner.Scan() {
			s := scanner.Text()
			if s != "" {
				Infof(context.TODO(), "ES日志: 请求参数: %s; 时间: %+v; 消耗：%+v", s, t, ti)
			}
		}
	}

	return nil
}

// RequestBodyEnabled 是否开启请求体日志
// return: 是否开启
func (logger EsLogger) RequestBodyEnabled() bool {
	return logger.RequestEnabled
}

// ResponseBodyEnabled 是否开启响应体日志
// return: 是否开启
func (logger EsLogger) ResponseBodyEnabled() bool {
	return logger.ResponseEnabled
}
