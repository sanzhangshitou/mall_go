package admin

import (
	"mall/app/Support/database"
	"mall/app/Support/response"

	"github.com/gin-gonic/gin"
)

// DashboardController 管理后台仪表盘
type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

type DashboardStats struct {
	TotalProducts   int64 `json:"total_products"`
	OnSaleProducts  int64 `json:"on_sale_products"`
	DraftProducts   int64 `json:"draft_products"`
	TotalCategories int64 `json:"total_categories"`
	TotalSkus       int64 `json:"total_skus"`
}

// Stats 获取仪表盘统计
func (ctrl *DashboardController) Stats(c *gin.Context) {
	var stats DashboardStats

	database.DB().Table("product").Where("is_deleted = ?", false).Count(&stats.TotalProducts)
	database.DB().Table("product").Where("is_deleted = ? AND status = ?", false, 1).Count(&stats.OnSaleProducts)
	database.DB().Table("product").Where("is_deleted = ? AND status = ?", false, 0).Count(&stats.DraftProducts)
	database.DB().Table("product_category").Where("is_deleted = ?", false).Count(&stats.TotalCategories)
	database.DB().Table("product_sku").Where("is_deleted = ?", false).Count(&stats.TotalSkus)

	response.Success(c, stats)
}
