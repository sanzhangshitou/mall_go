package middleware

import (
	"time"

	"mall/app/Support/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLog 请求日志中间件
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("请求日志",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		)
	}
}
