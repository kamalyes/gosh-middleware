/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 10:16:28
 * @FilePath: \gosh-middleware\trace.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"context"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/gosh"
	"github.com/kamalyes/gosh/constants"
)

// NewTraceIDContext 将追踪ID存储到上下文中
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, constants.TraceIdKey, traceID)
}

// GetTraceID 从上下文中获取追踪ID
func GetTraceID(ctx context.Context) (string, bool) {
	if v := ctx.Value(constants.TraceIdKey); v != nil {
		if id, ok := v.(string); ok && id != "" {
			return id, true
		}
	}
	return "", false
}

// TraceMiddleware 追踪中间件处理函数
func TraceMiddleware(skippers ...SkipperFunc) gosh.HandlerFunc {
	return func(c *gosh.Context) error {
		// 检查是否跳过中间件
		if SkipHandler(c, skippers...) {
			c.Next()
			return nil
		}

		// 从请求头中获取追踪ID，若为空则生成一个新的追踪ID
		traceID := c.Header(constants.TraceIdKey)
		if traceID == "" {
			traceID = osx.HashUnixMicroCipherText()
		}

		// 将追踪ID存储到上下文中，并设置响应头中的追踪ID
		ctx := NewTraceIDContext(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)
		c.ResponseWriter.Header().Set(constants.TraceIdKey, traceID)
		c.Next()
		return nil
	}
}
