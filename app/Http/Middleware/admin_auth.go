package middleware

import (
	"github.com/gin-gonic/gin"
)

// AdminAuth 管理后台认证中间件（占位，后续接入 JWT/Session）
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 从 Header 或 Cookie 中校验管理员 token
		// token := c.GetHeader("Authorization")
		// if token == "" || !validateAdminToken(token) {
		//     response.Forbidden(c, "无权限访问")
		//     c.Abort()
		//     return
		// }
		c.Next()
	}
}

// validateAdminToken 校验 token 的占位函数
func validateAdminToken(_ string) bool {
	// TODO: JWT 解析 + Redis 校验
	return true
}
