package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSON 统一 JSON 响应结构
type JSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// ---------- 成功响应 ----------

// Success 成功 (data可为空)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, JSON{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

// SuccessWithMsg 成功 + 自定义消息
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, JSON{
		Code:    0,
		Message: msg,
		Data:    data,
	})
}

// Page 分页成功
func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, JSON{
		Code:    0,
		Message: "ok",
		Data: PageData{
			List:       list,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}

// Created 创建成功
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, JSON{
		Code:    0,
		Message: "创建成功",
		Data:    data,
	})
}

// Updated 更新成功
func Updated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, JSON{
		Code:    0,
		Message: "更新成功",
		Data:    data,
	})
}

// Deleted 删除成功
func Deleted(c *gin.Context) {
	c.JSON(http.StatusOK, JSON{
		Code:    0,
		Message: "删除成功",
	})
}

// ---------- 失败响应 ----------

// Error 通用错误
func Error(c *gin.Context, httpCode int, code int, msg string) {
	c.JSON(httpCode, JSON{
		Code:    code,
		Message: msg,
	})
}

// BadRequest 参数错误
func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, 400, msg)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, 404, msg)
}

// InternalError 服务器内部错误
func InternalError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, 500, msg)
}

// ValidationError 表单验证失败
func ValidationError(c *gin.Context, msg string) {
	Error(c, http.StatusUnprocessableEntity, 422, msg)
}

// Forbidden 无权限
func Forbidden(c *gin.Context, msg string) {
	Error(c, http.StatusForbidden, 403, msg)
}
