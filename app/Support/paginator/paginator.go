package paginator

import (
	"strconv"

	"mall/config"

	"github.com/gin-gonic/gin"
)

// Paginator 分页参数
type Paginator struct {
	Page     int
	PageSize int
	Offset   int
}

// Parse 从请求中解析分页参数
func Parse(c *gin.Context) Paginator {
	cfg := config.Page()
	page := 1
	pageSize := cfg.DefaultSize

	if p := c.Query("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if n, err := strconv.Atoi(ps); err == nil && n > 0 {
			if n > cfg.MaxSize {
				n = cfg.MaxSize
			}
			pageSize = n
		}
	}

	return Paginator{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}
