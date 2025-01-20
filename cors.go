/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 10:55:15
 * @FilePath: \gosh-middleware\cors.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/kamalyes/go-config/pkg/cors"
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/gosh"
	"github.com/kamalyes/gosh/constants"
)

// NewCorsMiddleware 跨域中间件
func NewCorsMiddleware(config cors.Cors) gosh.HandlerFunc {
	return func(c *gosh.Context) error {
		origin := c.Header(constants.HeaderOriginKey)

		if !config.AllowedAllOrigins && !isOriginAllowed(origin, config.AllowedOrigins) {
			c.AbortWithStatus(config.OptionsResponseCode)
			return nil
		}

		setCorsHeaders(c, origin, config)
		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(config.OptionsResponseCode)
			return nil
		}

		defer recoverFromPanic()
		c.Next()
		return nil
	}
}

// setCorsHeaders 设置Cors头部信息
func setCorsHeaders(c *gosh.Context, origin string, config cors.Cors) {
	// 接收客户端发送的origin （重要）
	c.SetHeader(constants.AccessControlAllowOriginKey, origin)
	// 服务器支持的所有跨域请求的方法
	c.SetHeader(constants.AccessControlAllowMethodsKey, strings.Join(config.AllowedMethods, ","))
	// 允许跨域设置可以返回其他子段，可以自定义字段
	c.SetHeader(constants.AccessControlAllowHeadersKey, strings.Join(config.AllowedHeaders, ","))
	// 允许浏览器（客户端）可以解析的头部 （重要）
	c.SetHeader(constants.AccessControlExposeHeadersKey, strings.Join(config.ExposedHeaders, ","))
	// 设置缓存时间
	c.SetHeader(constants.AccessControlMaxAgeKey, config.MaxAge)
	// 允许客户端传递校验信息比如 cookie (重要)
	c.SetHeader(constants.AccessControlAllowCredentialsKey, convert.MustString(config.AllowCredentials))
}

// isOriginAllowed 检查请求的 Origin 是否在允许的列表中
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if len(allowedOrigins) == 0 {
		return false
	}
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || strings.TrimSpace(allowedOrigin) == origin {
			return true
		}
	}
	return false
}

func recoverFromPanic() {
	if err := recover(); err != nil {
		log.Printf("Panic info: %v", err)
	}
}
