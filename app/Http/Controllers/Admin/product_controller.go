package admin

import (
	"encoding/json"
	"strconv"

	adminreq "mall/app/Http/Requests/Admin"
	adminres "mall/app/Http/Resources/Admin"
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
	return &ProductController{service: services.NewProductService()}
}

// Index 商品列表（可选含已删除）
func (ctrl *ProductController) Index(c *gin.Context) {
	var req adminreq.ProductListRequest
	_ = c.ShouldBindQuery(&req)
	pg := paginator.Parse(c)

	query := services.ProductListQuery{
		CategoryID: req.CategoryID, Status: req.Status,
		Keyword: req.Keyword, SortField: req.SortField, SortOrder: req.SortOrder,
		Page: pg.Page, PageSize: pg.PageSize, Offset: pg.Offset,
	}

	products, total, err := ctrl.service.AdminList(query, req.IncludeDeleted)
	if err != nil {
		logger.Error("管理后台商品列表失败", zap.Error(err))
		response.InternalError(c, "获取失败")
		return
	}
	list := adminres.ToProductList(products)

	// 批量补分类名
	if len(products) > 0 {
		catSvc := services.NewCategoryService()
		ids := make([]int64, len(products))
		for i, p := range products {
			ids[i] = p.CategoryID
		}
		cats, _ := catSvc.GetByIDs(ids)
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

// Show 商品详情（含已删除）
func (ctrl *ProductController) Show(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	product, err := ctrl.service.AdminGetByID(id)
	if err != nil {
		response.InternalError(c, "获取失败")
		return
	}
	if product == nil {
		response.NotFound(c, "商品不存在")
		return
	}
	detail := adminres.ToProductDetail(product)
	if cat, _ := services.NewCategoryService().GetByID(product.CategoryID); cat != nil {
		detail.CategoryName = cat.Name
	}
	response.Success(c, detail)
}

// Store 创建商品
func (ctrl *ProductController) Store(c *gin.Context) {
	var req adminreq.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	product := &models.Product{
		CategoryID: req.CategoryID, Name: req.Name,
		Status: models.ProductStatus(req.Status), SortOrder: req.SortOrder,
	}
	if req.Images != nil {
		b, _ := json.Marshal(req.Images)
		product.Images = b
	}
	if req.Specs != nil {
		b, _ := json.Marshal(req.Specs)
		product.Specs = b
	}
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
		response.InternalError(c, "创建失败")
		return
	}
	response.Created(c, adminres.ToProductResource(product))
}

// Update 更新商品
func (ctrl *ProductController) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	product, _ := ctrl.service.AdminGetByID(id)
	if product == nil {
		response.NotFound(c, "商品不存在")
		return
	}
	var req adminreq.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
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
		response.InternalError(c, "更新失败")
		return
	}
	response.Updated(c, adminres.ToProductResource(product))
}

// Destroy 软删除
func (ctrl *ProductController) Destroy(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := ctrl.service.Delete(id); err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Deleted(c)
}

// Restore 恢复
func (ctrl *ProductController) Restore(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := ctrl.service.Restore(id); err != nil {
		response.InternalError(c, "恢复失败")
		return
	}
	response.SuccessWithMsg(c, "恢复成功", nil)
}

// UpdateStatus 修改商品状态
func (ctrl *ProductController) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req adminreq.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	if err := ctrl.service.UpdateStatus(id, models.ProductStatus(req.Status)); err != nil {
		response.InternalError(c, "状态更新失败")
		return
	}
	product, _ := ctrl.service.AdminGetByID(id)
	response.Updated(c, adminres.ToProductResource(product))
}
