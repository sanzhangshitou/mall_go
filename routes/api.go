package routes

import (
	api "mall/app/Http/Controllers/Api"
	middleware "mall/app/Http/Middleware"

	"github.com/gin-gonic/gin"
)

// Register 注册所有 API 路由
func Register(router *gin.Engine) {
	// 全局中间件
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLog())
	router.Use(middleware.CORS())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 商品分类
		categoryCtrl := api.NewCategoryController()
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryCtrl.Index)
			categories.GET("/:id", categoryCtrl.Show)
			categories.POST("", categoryCtrl.Store)
			categories.PUT("/:id", categoryCtrl.Update)
			categories.DELETE("/:id", categoryCtrl.Destroy)
		}

		// 商品
		productCtrl := api.NewProductController()
		products := v1.Group("/products")
		{
			products.GET("", productCtrl.Index)
			products.GET("/:id", productCtrl.Show)
			products.POST("", productCtrl.Store)
			products.PUT("/:id", productCtrl.Update)
			products.DELETE("/:id", productCtrl.Destroy)
		}

		// 商品 SKU (嵌套路由)
		skuCtrl := api.NewProductSkuController()
		skus := v1.Group("/products/:id/skus")
		{
			skus.GET("", skuCtrl.Index)
			skus.GET("/:id", skuCtrl.Show)
			skus.POST("", skuCtrl.Store)
			skus.PUT("/:id", skuCtrl.Update)
			skus.DELETE("/:id", skuCtrl.Destroy)
		}
	}
}
