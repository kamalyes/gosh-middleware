/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 15:15:17
 * @FilePath: \gosh-middleware\body.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/kamalyes/gosh"
	"github.com/kamalyes/gosh/constants"
)

var (
	maxMemory int64 = 64 << 20 // 64 MB
	mu        sync.Mutex
)

// GetMaxMemory 获取最大内存
func GetMaxMemory() int64 {
	return syncx.WithLockReturnValue(&mu, func() int64 {
		return maxMemory
	})
}

// SetMaxMemory 设置最大内存
func SetMaxMemory(memory int64) {
	syncx.WithLock(&mu, func() {
		maxMemory = memory
	})
}

func CopyBodyMiddleware(skippers ...SkipperFunc) gosh.HandlerFunc {
	return func(c *gosh.Context) error {
		if SkipHandler(c, skippers...) || c.Request.Body == nil {
			c.Next()
			return nil
		}

		var requestBody []byte
		isGzip := false
		safe := &io.LimitedReader{R: c.Request.Body, N: GetMaxMemory()}

		// 检查请求是否使用gzip压缩
		if c.Header(constants.HeaderContentEncodingKey) == constants.ContentEncodingGzip {
			reader, err := gzip.NewReader(safe)
			if err == nil {
				defer reader.Close() // 确保在使用后关闭 reader
				isGzip = true
				requestBody, _ = io.ReadAll(reader)
			}
		}

		// 如果不是gzip压缩或解压缩时出错，则直接读取请求体
		if !isGzip {
			requestBody, _ = io.ReadAll(safe)
		}

		// 关闭原始请求体，创建新的 MaxBytesReader，将复制的请求体作为缓冲区
		c.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		c.Request.Body = http.MaxBytesReader(c.ResponseWriter, io.NopCloser(bf), GetMaxMemory())
		c.SetContextValue(constants.RequestBody, requestBody)
		c.Next()
		return nil
	}
}
