package routes

import (
	adminctrl "mall/app/Http/Controllers/Admin"
	middleware "mall/app/Http/Middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAdmin 注册管理后台路由 /admin
func RegisterAdmin(router *gin.Engine) {
	admin := router.Group("/admin")

	// 登录接口 — 不受 AdminAuth 限制
	authCtrl := adminctrl.NewAuthController()
	admin.POST("/login", authCtrl.Login)

	// 其余接口需认证
	admin.Use(middleware.AdminAuth())
	{
		// Dashboard
		dashboardCtrl := adminctrl.NewDashboardController()
		admin.GET("/dashboard", dashboardCtrl.Stats)

		// 分类
		catCtrl := adminctrl.NewCategoryController()
		categories := admin.Group("/categories")
		{
			categories.GET("", catCtrl.Index)
			categories.GET("/:id", catCtrl.Show)
			categories.POST("", catCtrl.Store)
			categories.PUT("/:id", catCtrl.Update)
			categories.DELETE("/:id", catCtrl.Destroy)
			categories.PATCH("/:id/restore", catCtrl.Restore)
		}

		// 商品
		productCtrl := adminctrl.NewProductController()
		products := admin.Group("/products")
		{
			products.GET("", productCtrl.Index)
			products.GET("/:id", productCtrl.Show)
			products.POST("", productCtrl.Store)
			products.PUT("/:id", productCtrl.Update)
			products.DELETE("/:id", productCtrl.Destroy)
			products.PATCH("/:id/restore", productCtrl.Restore)
			products.PATCH("/:id/status", productCtrl.UpdateStatus)
		}

		// 商品 SKU
		skuCtrl := adminctrl.NewSkuController()
		skus := admin.Group("/products/:id/skus")
		{
			skus.GET("", skuCtrl.Index)
			skus.GET("/:skuId", skuCtrl.Show)
			skus.POST("", skuCtrl.Store)
			skus.PUT("/:skuId", skuCtrl.Update)
			skus.DELETE("/:skuId", skuCtrl.Destroy)
			skus.PATCH("/:skuId/restore", skuCtrl.Restore)
		}
	}
}
