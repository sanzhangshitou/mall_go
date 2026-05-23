package api

import (
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

type CategoryController struct {
	service *services.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		service: services.NewCategoryService(),
	}
}

// Index 分类列表 (树形)
func (ctrl *CategoryController) Index(c *gin.Context) {
	categories, err := ctrl.service.GetAll()
	if err != nil {
		logger.Error("获取分类列表失败", zap.Error(err))
		response.InternalError(c, "获取分类列表失败")
		return
	}

	tree := apires.ToCategoryTree(categories)
	response.Success(c, tree)
}

// Show 分类详情
func (ctrl *CategoryController) Show(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	cat, err := ctrl.service.GetByID(id)
	if err != nil {
		logger.Error("获取分类详情失败", zap.Error(err))
		response.InternalError(c, "获取分类详情失败")
		return
	}
	if cat == nil {
		response.NotFound(c, "分类不存在")
		return
	}

	response.Success(c, apires.ToCategoryResource(cat))
}

// Store 创建分类
func (ctrl *CategoryController) Store(c *gin.Context) {
	var req apireq.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "参数校验失败: "+err.Error())
		return
	}

	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}

	cat := &models.ProductCategory{
		Name:      req.Name,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		IsEnabled: isEnabled,
	}

	if err := ctrl.service.Create(cat); err != nil {
		logger.Error("创建分类失败", zap.Error(err))
		response.InternalError(c, "创建分类失败")
		return
	}

	response.Created(c, apires.ToCategoryResource(cat))
}

// Update 更新分类
func (ctrl *CategoryController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	cat, err := ctrl.service.GetByID(id)
	if err != nil {
		logger.Error("获取分类失败", zap.Error(err))
		response.InternalError(c, "获取分类失败")
		return
	}
	if cat == nil {
		response.NotFound(c, "分类不存在")
		return
	}

	var req apireq.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "参数校验失败: "+err.Error())
		return
	}

	if req.Name != "" {
		cat.Name = req.Name
	}
	if req.ParentID != nil {
		cat.ParentID = *req.ParentID
	}
	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}
	if req.IsEnabled != nil {
		cat.IsEnabled = *req.IsEnabled
	}

	if err := ctrl.service.Update(cat); err != nil {
		logger.Error("更新分类失败", zap.Error(err))
		response.InternalError(c, "更新分类失败")
		return
	}

	response.Updated(c, apires.ToCategoryResource(cat))
}

// Destroy 删除分类 (软删除)
func (ctrl *CategoryController) Destroy(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	if err := ctrl.service.Delete(id); err != nil {
		logger.Error("删除分类失败", zap.Error(err))
		response.InternalError(c, "删除分类失败")
		return
	}

	response.Deleted(c)
}
