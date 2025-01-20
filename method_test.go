/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 10:16:28
 * @FilePath: \gosh-middleware\method_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"net/http"
	"testing"

	"github.com/kamalyes/gosh"
	"github.com/stretchr/testify/assert"
)

// 创建一个测试上下文
func createTestMethodContext(method, path string) *gosh.Context {
	req, _ := http.NewRequest(method, path, nil)
	return &gosh.Context{
		Request: req,
	}
}

func TestAllowPathPrefixSkipper(t *testing.T) {
	skipper := AllowPathPrefixSkipper("/api", "/health")

	tests := []struct {
		path     string
		expected bool
	}{
		{"/api/v1/resource", true},
		{"/health", true},
		{"/other", false},
	}

	for _, test := range tests {
		ctx := createTestMethodContext("GET", test.path)
		result := skipper(ctx)
		assert.Equal(t, test.expected, result, "Expected result for path %s to be %v", test.path, test.expected)
	}
}

func TestAllowMethodAndPathPrefixSkipper(t *testing.T) {
	skipper := AllowMethodAndPathPrefixSkipper("GET/api/v1", "POST/api/v2")

	tests := []struct {
		method   string
		path     string
		expected bool
	}{
		{"GET", "/api/v1/resource", true},
		{"POST", "/api/v2/resource", true},
		{"GET", "/api/v2/resource", false},
		{"DELETE", "/api/v1/resource", false},
	}

	for _, test := range tests {
		ctx := createTestMethodContext(test.method, test.path)
		result := skipper(ctx)
		assert.Equal(t, test.expected, result, "Expected result for method %s and path %s to be %v", test.method, test.path, test.expected)
	}
}

func TestJoinRouter(t *testing.T) {
	tests := []struct {
		method   string
		path     string
		expected string
	}{
		{"get", "api/v1/resource", "GET/api/v1/resource"},
		{"POST", "/api/v2/resource", "POST/api/v2/resource"},
		{"PUT", "api/v3/resource", "PUT/api/v3/resource"},
	}

	for _, test := range tests {
		result := JoinRouter(test.method, test.path)
		assert.Equal(t, test.expected, result, "Expected JoinRouter(%s, %s) to be %s", test.method, test.path, test.expected)
	}
}

func TestSkipHandler(t *testing.T) {
	skipper1 := AllowPathPrefixSkipper("/api")
	skipper2 := AllowPathPrefixSkipper("/health")

	tests := []struct {
		path     string
		expected bool
	}{
		{"/api/v1/resource", true},
		{"/health", true},
		{"/other", false},
	}

	for _, test := range tests {
		ctx := createTestMethodContext("GET", test.path)
		result := SkipHandler(ctx, skipper1, skipper2)
		assert.Equal(t, test.expected, result, "Expected SkipHandler(%s) to be %v", test.path, test.expected)
	}
}
