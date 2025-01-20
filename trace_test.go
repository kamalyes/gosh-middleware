/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 10:16:28
 * @FilePath: \gosh-middleware\trace_test.go
 * @Description: 对 trace 包进行单元测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamalyes/gosh"
	"github.com/kamalyes/gosh/constants"
	"github.com/stretchr/testify/assert"
)

func TestNewTraceIDContext(t *testing.T) {
	// 测试将追踪ID存储到上下文中
	traceID := "123456"
	ctx := NewTraceIDContext(context.Background(), traceID)

	retrievedID, ok := GetTraceID(ctx)
	assert.True(t, ok, "Expected to retrieve a trace ID from context.")
	assert.Equal(t, traceID, retrievedID, "Expected trace ID should match the stored value.")
}

func TestGetTraceID_NoValue(t *testing.T) {
	// 测试从上下文中获取追踪ID时没有值的情况
	ctx := context.Background()
	retrievedID, ok := GetTraceID(ctx)

	assert.False(t, ok, "Expected to not retrieve a trace ID from context.")
	assert.Empty(t, retrievedID, "Expected trace ID should be empty.")
}

func TestTraceMiddleware(t *testing.T) {
	// 创建一个模拟的上下文
	c := &gosh.Context{
		Request: &http.Request{
			Header: make(http.Header),
		},
		ResponseWriter: httptest.NewRecorder(),
	}

	// 测试中间件，确保它可以获取追踪ID并设置响应头
	request := TraceMiddleware()
	err := request(c)

	assert.NoError(t, err, "Expected no error from TraceMiddleware.")
	assert.NotEmpty(t, c.ResponseWriter.Header().Get(constants.TraceIdKey), "Expected TraceRequestKey in response headers.")
}

func TestTraceMiddleware_Skip(t *testing.T) {
	// 创建一个模拟的上下文
	c := &gosh.Context{
		Request: &http.Request{
			Header: make(http.Header),
		},
		ResponseWriter: httptest.NewRecorder(),
	}

	// 测试跳过中间件的情况
	request := TraceMiddleware(func(c *gosh.Context) bool {
		return true // 总是跳过
	})

	err := request(c)

	assert.NoError(t, err, "Expected no error from TraceMiddleware when skipped.")
	assert.Empty(t, c.ResponseWriter.Header().Get(constants.TraceIdKey), "Expected no TraceRequestKey in response headers when skipped.")
}
