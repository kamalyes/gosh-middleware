/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 13:15:55
 * @FilePath: \gosh-middleware\cors_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamalyes/go-config/pkg/cors"
	"github.com/kamalyes/gosh"
	"github.com/stretchr/testify/assert"
)

func TestNewCorsMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		origin          string
		corsConfig      cors.Cors
		expectedStatus  int
		expectedHeaders map[string]string
	}{
		{
			name:   "Allowed Origin",
			origin: "http://allowed-origin.com",
			corsConfig: cors.Cors{
				AllowedAllOrigins:   false,
				AllowedOrigins:      []string{"http://allowed-origin.com"},
				AllowedMethods:      []string{"GET", "POST"},
				AllowedHeaders:      []string{"Content-Type"},
				ExposedHeaders:      []string{"X-Custom-Header"},
				MaxAge:              "3600",
				AllowCredentials:    true,
				OptionsResponseCode: http.StatusOK,
			},
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://allowed-origin.com",
				"Access-Control-Allow-Methods":     "GET,POST",
				"Access-Control-Allow-Headers":     "Content-Type",
				"Access-Control-Expose-Headers":    "X-Custom-Header",
				"Access-Control-Max-Age":           "3600",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:   "Disallowed Origin",
			origin: "http://disallowed-origin.com",
			corsConfig: cors.Cors{
				AllowedAllOrigins:   false,
				AllowedOrigins:      []string{"http://allowed-origin.com"},
				OptionsResponseCode: http.StatusForbidden,
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "Wildcard Origin",
			origin: "http://any-origin.com",
			corsConfig: cors.Cors{
				AllowedAllOrigins:   true,
				OptionsResponseCode: http.StatusOK,
			},
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "http://any-origin.com",
			},
		},
		{
			name:   "Preflight Request",
			origin: "http://allowed-origin.com",
			corsConfig: cors.Cors{
				AllowedAllOrigins:   false,
				AllowedOrigins:      []string{"http://allowed-origin.com"},
				AllowedMethods:      []string{"GET", "POST"},
				OptionsResponseCode: http.StatusNoContent,
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gosh context
			c := &gosh.Context{
				Request:        httptest.NewRequest(http.MethodOptions, "/", nil),
				ResponseWriter: httptest.NewRecorder(),
			}
			c.Request.Header.Set("Origin", tt.origin)

			// Create the middleware
			middleware := NewCorsMiddleware(tt.corsConfig)

			// Call the middleware
			err := middleware(c)
			if err != nil {
				t.Fatalf("Middleware returned an error: %v", err)
			}

			// Check the response status
			assert.Equal(t, tt.expectedStatus, c.Status)

			// Check headers if expected
			for key, value := range tt.expectedHeaders {
				assert.Equal(t, value, c.ResponseWriter.Header().Get(key))
			}
		})
	}
}
