package api

import (
	"encoding/json"
	"strconv"

	apireq "mall/app/Http/Requests/Api"
	apires "mall/app/Http/Resources/Api"
	models "mall/app/Models"
	services "mall/app/Services"
	"mall/app/Support/logger"
	"mall/app/Support/paginator"
	"mall/app/Support/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductController struct {
	service *services.ProductService
}

func NewProductController() *ProductController {
	return &ProductController{
		service: services.NewProductService(),
	}
}

// Index 商品列表
func (ctrl *ProductController) Index(c *gin.Context) {
	var req apireq.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, "参数校验失败: "+err.Error())
		return
	}

	pg := paginator.Parse(c)

	query := services.ProductListQuery{
		CategoryID: req.CategoryID,
		Status:     req.Status,
		Keyword:    req.Keyword,
		SortField:  req.SortField,
		SortOrder:  req.SortOrder,
		Page:       pg.Page,
		PageSize:   pg.PageSize,
		Offset:     pg.Offset,
	}

	products, total, err := ctrl.service.List(query)
	if err != nil {
		logger.Error("获取商品列表失败", zap.Error(err))
		response.InternalError(c, "获取商品列表失败")
		return
	}

	// 转换为资源
	list := apires.ToProductList(products)

	// 批量补分类名称 (符合规范: 不连表, 查完后补数据)
	if len(products) > 0 {
		categorySvc := services.NewCategoryService()
		categoryIDs := make([]int64, 0, len(products))
		for _, p := range products {
			categoryIDs = append(categoryIDs, p.CategoryID)
		}
		cats, _ := categorySvc.GetByIDs(categoryIDs)
		catMap := make(map[int64]string, len(cats))
		for _, cat := range cats {
			catMap[cat.ID] = cat.Name
		}
		for _, r := range list {
			if name, ok := catMap[r.CategoryID]; ok {
				r.CategoryName = name
			}
		}
	}

	response.Page(c, list, total, pg.Page, pg.PageSize)
}

// Show 商品详情
func (ctrl *ProductController) Show(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的商品ID")
		return
	}

	product, err := ctrl.service.GetByID(id)
	if err != nil {
		logger.Error("获取商品详情失败", zap.Error(err))
		response.InternalError(c, "获取商品详情失败")
		return
	}
	if product == nil {
		response.NotFound(c, "商品不存在")
		return
	}

	detail := apires.ToProductDetail(product)

	// 补分类名称
	categorySvc := services.NewCategoryService()
	if cat, _ := categorySvc.GetByID(product.CategoryID); cat != nil {
		detail.CategoryName = cat.Name
	}

	response.Success(c, detail)
}

// Store 创建商品
func (ctrl *ProductController) Store(c *gin.Context) {
	var req apireq.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "参数校验失败: "+err.Error())
		return
	}

	// 序列化 Images
	var imagesJSON json.RawMessage
	if req.Images != nil {
		b, _ := json.Marshal(req.Images)
		imagesJSON = b
	}

	// 序列化 Specs
	var specsJSON json.RawMessage
	if req.Specs != nil {
		b, _ := json.Marshal(req.Specs)
		specsJSON = b
	}

	product := &models.Product{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Images:     imagesJSON,
		Status:     models.ProductStatus(req.Status),
		Specs:      specsJSON,
		SortOrder:  req.SortOrder,
	}

	// 可选字段
	if req.Subtitle != "" {
		product.Subtitle.Scan(req.Subtitle)
	}
	if req.Description != "" {
		product.Description.Scan(req.Description)
	}
	if req.MainImage != "" {
		product.MainImage.Scan(req.MainImage)
	}
	if req.Unit != "" {
		product.Unit.Scan(req.Unit)
	}

	if err := ctrl.service.Create(product); err != nil {
		logger.Error("创建商品失败", zap.Error(err))
		response.InternalError(c, "创建商品失败")
		return
	}

	response.Created(c, apires.ToProductResource(product))
}

// Update 更新商品
func (ctrl *ProductController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的商品ID")
		return
	}

	product, err := ctrl.service.GetByID(id)
	if err != nil {
		logger.Error("获取商品失败", zap.Error(err))
		response.InternalError(c, "获取商品失败")
		return
	}
	if product == nil {
		response.NotFound(c, "商品不存在")
		return
	}

	var req apireq.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "参数校验失败: "+err.Error())
		return
	}

	// 按需更新字段
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Subtitle != nil {
		product.Subtitle.Scan(*req.Subtitle)
	}
	if req.Description != nil {
		product.Description.Scan(*req.Description)
	}
	if req.MainImage != nil {
		product.MainImage.Scan(*req.MainImage)
	}
	if req.Images != nil {
		b, _ := json.Marshal(*req.Images)
		product.Images = b
	}
	if req.Status != nil {
		product.Status = models.ProductStatus(*req.Status)
	}
	if req.Specs != nil {
		b, _ := json.Marshal(*req.Specs)
		product.Specs = b
	}
	if req.Unit != nil {
		product.Unit.Scan(*req.Unit)
	}
	if req.SortOrder != nil {
		product.SortOrder = *req.SortOrder
	}

	if err := ctrl.service.Update(product); err != nil {
		logger.Error("更新商品失败", zap.Error(err))
		response.InternalError(c, "更新商品失败")
		return
	}

	// 重新获取完整详情
	updated, _ := ctrl.service.GetByID(product.ID)
	response.Updated(c, apires.ToProductResource(updated))
}

// Destroy 删除商品 (软删除)
func (ctrl *ProductController) Destroy(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的商品ID")
		return
	}

	if err := ctrl.service.Delete(id); err != nil {
		logger.Error("删除商品失败", zap.Error(err))
		response.InternalError(c, "删除商品失败")
		return
	}

	response.Deleted(c)
}
