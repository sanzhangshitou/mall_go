package middleware

import (
	"mall/app/Support/logger"
	"mall/app/Support/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 异常恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic 恢复", zap.Any("error", err))

				// 判断是否是 gin 的 error 类型
				if e, ok := err.(error); ok {
					response.InternalError(c, "服务器内部错误: "+e.Error())
				} else {
					response.InternalError(c, "服务器内部错误")
				}

				c.Abort()
			}
		}()
		c.Next()
	}
}
