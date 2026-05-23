package admin

import (
	adminreq "mall/app/Http/Requests/Admin"
	"mall/app/Support/logger"
	"mall/app/Support/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthController 管理后台认证
type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

// Login 模拟登录
func (ctrl *AuthController) Login(c *gin.Context) {
	var req adminreq.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 模拟认证校验
	if req.Username != "admin" || req.Password != "admin123" {
		response.Error(c, 401, 401, "用户名或密码错误")
		return
	}

	logger.Info("管理员登录成功", zap.String("username", req.Username))

	response.Success(c, LoginResponse{
		Token:    "mock-token-admin-" + req.Username,
		Username: req.Username,
	})
}
