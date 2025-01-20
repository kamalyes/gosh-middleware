/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 15:15:15
 * @FilePath: \gosh-middleware\pprof_test.go
 * @Description: 对 pprof 包进行单元测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/kamalyes/gosh"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// 创建一个 HTTP 测试请求和响应
	req := httptest.NewRequest("GET", "/debug/pprof/sysinfo", nil)
	w := httptest.NewRecorder()

	// 创建一个新的上下文
	ctx := &gosh.Context{
		Request:        req,
		ResponseWriter: w,
	}

	// 调用 Handler
	err := PprofHandler(ctx)
	assert.NoError(t, err, "Handler should not return an error")

	// 检查响应状态码
	assert.Equal(t, 200, w.Code, "Expected status code 200")

	// 检查响应内容是否包含关键字（例如服务器名称）
	assert.Contains(t, w.Body.String(), "服务器", "Response should contain '服务器'")
	assert.Contains(t, w.Body.String(), "运行时间", "Response should contain '运行时间'")
	assert.Contains(t, w.Body.String(), "goroutine数量", "Response should contain 'goroutine数量'")
	assert.Contains(t, w.Body.String(), "CPU核数", "Response should contain 'CPU核数'")
}
