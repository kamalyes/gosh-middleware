/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 17:25:16
 * @FilePath: \gosh-middleware\sysinfo_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSystemInfoMonitor 持续监听 SystemInfo 的变化
func TestSystemInfoMonitor(t *testing.T) {
	// 设置检查间隔时间
	checkInterval := time.Second * 5
	// 使用 goroutine 来定期检查 SystemInfo 的变化
	go func() {
		for {
			time.Sleep(checkInterval)
			sysInfo := NewSystemInfo()
			assert.NotEmpty(t, sysInfo.ServerName, "ServerName should not be empty")
			assert.NotEmpty(t, sysInfo.Runtime, "Runtime should not be empty")
			assert.NotEmpty(t, sysInfo.GoroutineNum, "GoroutineNum should not be empty")
			assert.NotEmpty(t, sysInfo.CPUNum, "CPUNum should not be empty")
			assert.NotEmpty(t, sysInfo.UsedMem, "UsedMem should not be empty")
			assert.NotEmpty(t, sysInfo.HeapInuse, "HeapInuse should not be empty")
			assert.NotEmpty(t, sysInfo.TotalMem, "TotalMem should not be empty")
			assert.NotEmpty(t, sysInfo.SysMem, "SysMem should not be empty")
			assert.NotEmpty(t, sysInfo.Lookups, "Lookups should not be empty")
			assert.NotEmpty(t, sysInfo.Mallocs, "Mallocs should not be empty")
			assert.NotEmpty(t, sysInfo.Frees, "Frees should not be empty")
			assert.NotEmpty(t, sysInfo.LastGCTime, "LastGCTime should not be empty")
			assert.NotEmpty(t, sysInfo.NextGC, "NextGC should not be empty")
			assert.NotEmpty(t, sysInfo.PauseTotalNs, "PauseTotalNs should not be empty")
			assert.NotEmpty(t, sysInfo.PauseNs, "PauseNs should not be empty")
		}
	}()
}
