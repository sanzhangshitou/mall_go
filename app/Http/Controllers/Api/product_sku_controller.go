package api

import (
	"encoding/json"
	"strconv"

	apireq "mall/app/Http/Requests/Api"
	apires "mall/app/Http/Resources/Api"
	models "mall/app/Models"
	services "mall/app/Services"
	"mall/app/Support/logger"
	"mall/app/Support/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductSkuController struct {
	service *services.ProductSkuService
}

func NewProductSkuController() *ProductSkuController {
	return &ProductSkuController{
		service: services.NewProductSkuService(),
	}
}

// Index 获取商品下所有 SKU
func (ctrl *ProductSkuController) Index(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的商品ID")
		return
	}

	skus, err := ctrl.service.GetByProductID(productID)
	if err != nil {
		logger.Error("获取SKU列表失败", zap.Error(err))
		response.InternalError(c, "获取SKU列表失败")
		return
	}

	response.Success(c, apires.ToSkuList(skus))
}

// Show SKU 详情
func (ctrl *ProductSkuController) Show(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的SKU ID")
		return
	}

	sku, err := ctrl.service.GetByID(id)
	if err != nil {
		logger.Error("获取SKU详情失败", zap.Error(err))
		response.InternalError(c, "获取SKU详情失败")
		return
	}
	if sku == nil {
		response.NotFound(c, "SKU不存在")
		return
	}

	response.Success(c, apires.ToSkuResource(sku))
}

// Store 创建SKU
func (ctrl *ProductSkuController) Store(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的商品ID")
		return
	}

	var req apireq.CreateSkuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "参数校验失败: "+err.Error())
		return
	}

	// 序列化 SpecValues
	var specJSON json.RawMessage
	if req.SpecValues != nil {
		b, _ := json.Marshal(req.SpecValues)
		specJSON = b
	}

	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}

	sku := &models.ProductSku{
		ProductID:  productID,
		SpecValues: specJSON,
		Price:      req.Price,
		Stock:      req.Stock,
		IsEnabled:  isEnabled,
		SortOrder:  req.SortOrder,
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
		logger.Error("创建SKU失败", zap.Error(err))
		response.InternalError(c, "创建SKU失败")
		return
	}

	response.Created(c, apires.ToSkuResource(sku))
}

// Update 更新SKU
func (ctrl *ProductSkuController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的SKU ID")
		return
	}

	sku, err := ctrl.service.GetByID(id)
	if err != nil {
		logger.Error("获取SKU失败", zap.Error(err))
		response.InternalError(c, "获取SKU失败")
		return
	}
	if sku == nil {
		response.NotFound(c, "SKU不存在")
		return
	}

	var req apireq.UpdateSkuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "参数校验失败: "+err.Error())
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
		logger.Error("更新SKU失败", zap.Error(err))
		response.InternalError(c, "更新SKU失败")
		return
	}

	response.Updated(c, apires.ToSkuResource(sku))
}

// Destroy 删除SKU (软删除)
func (ctrl *ProductSkuController) Destroy(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的SKU ID")
		return
	}

	if err := ctrl.service.Delete(id); err != nil {
		logger.Error("删除SKU失败", zap.Error(err))
		response.InternalError(c, "删除SKU失败")
		return
	}

	response.Deleted(c)
}
