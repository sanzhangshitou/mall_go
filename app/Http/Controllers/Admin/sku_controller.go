package admin

import (
	"encoding/json"
	"strconv"

	adminreq "mall/app/Http/Requests/Admin"
	adminres "mall/app/Http/Resources/Admin"
	models "mall/app/Models"
	services "mall/app/Services"
	"mall/app/Support/logger"
	"mall/app/Support/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SkuController struct {
	service *services.ProductSkuService
}

func NewSkuController() *SkuController {
	return &SkuController{service: services.NewProductSkuService()}
}

// Index SKU 列表（含已删除）
func (ctrl *SkuController) Index(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	skus, err := ctrl.service.AdminGetByProductID(productID)
	if err != nil {
		logger.Error("管理后台SKU列表失败", zap.Error(err))
		response.InternalError(c, "获取失败")
		return
	}
	response.Success(c, adminres.ToSkuList(skus))
}

// Show SKU 详情（含已删除）
func (ctrl *SkuController) Show(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("skuId"), 10, 64)
	sku, err := ctrl.service.AdminGetByID(id)
	if err != nil {
		response.InternalError(c, "获取失败")
		return
	}
	if sku == nil {
		response.NotFound(c, "SKU不存在")
		return
	}
	response.Success(c, adminres.ToSkuResource(sku))
}

// Store 创建SKU
func (ctrl *SkuController) Store(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req adminreq.CreateSkuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}
	sku := &models.ProductSku{
		ProductID: productID, Price: req.Price, Stock: req.Stock,
		IsEnabled: isEnabled, SortOrder: req.SortOrder,
	}
	if req.SpecValues != nil {
		b, _ := json.Marshal(req.SpecValues)
		sku.SpecValues = b
	}
	if req.SkuCode != "" {
		sku.SkuCode.Scan(req.SkuCode)
	}
	if req.MarketPrice != nil {
		sku.MarketPrice.Scan(*req.MarketPrice)
	}
	if req.Image != "" {
		sku.Image.Scan(req.Image)
	}
	if err := ctrl.service.Create(sku); err != nil {
		response.InternalError(c, "创建失败")
		return
	}
	response.Created(c, adminres.ToSkuResource(sku))
}

// Update 更新SKU
func (ctrl *SkuController) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("skuId"), 10, 64)
	sku, _ := ctrl.service.AdminGetByID(id)
	if sku == nil {
		response.NotFound(c, "SKU不存在")
		return
	}
	var req adminreq.UpdateSkuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	if req.SkuCode != nil {
		sku.SkuCode.Scan(*req.SkuCode)
	}
	if req.SpecValues != nil {
		b, _ := json.Marshal(*req.SpecValues)
		sku.SpecValues = b
	}
	if req.Price != nil {
		sku.Price = *req.Price
	}
	if req.MarketPrice != nil {
		sku.MarketPrice.Scan(*req.MarketPrice)
	}
	if req.Stock != nil {
		sku.Stock = *req.Stock
	}
	if req.Image != nil {
		sku.Image.Scan(*req.Image)
	}
	if req.IsEnabled != nil {
		sku.IsEnabled = *req.IsEnabled
	}
	if req.SortOrder != nil {
		sku.SortOrder = *req.SortOrder
	}
	if err := ctrl.service.Update(sku); err != nil {
		response.InternalError(c, "更新失败")
		return
	}
	response.Updated(c, adminres.ToSkuResource(sku))
}

// Destroy 软删除
func (ctrl *SkuController) Destroy(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("skuId"), 10, 64)
	if err := ctrl.service.Delete(id); err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Deleted(c)
}

// Restore 恢复
func (ctrl *SkuController) Restore(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("skuId"), 10, 64)
	if err := ctrl.service.Restore(id); err != nil {
		response.InternalError(c, "恢复失败")
		return
	}
	response.SuccessWithMsg(c, "恢复成功", nil)
}
