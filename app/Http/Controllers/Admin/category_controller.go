package admin

import (
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

type CategoryController struct {
	service *services.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{service: services.NewCategoryService()}
}

// Index 分类列表（可选显示禁用/已删除）
func (ctrl *CategoryController) Index(c *gin.Context) {
	var req adminreq.CategoryListRequest
	_ = c.ShouldBindQuery(&req)

	categories, err := ctrl.service.AdminGetAll(req.IncludeDisabled, req.IncludeDeleted)
	if err != nil {
		logger.Error("管理后台获取分类列表失败", zap.Error(err))
		response.InternalError(c, "获取失败")
		return
	}
	tree := adminres.ToCategoryTree(categories)
	response.Success(c, tree)
}

// Show 分类详情（含已删除）
func (ctrl *CategoryController) Show(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cat, err := ctrl.service.AdminGetByID(id)
	if err != nil {
		response.InternalError(c, "获取失败")
		return
	}
	if cat == nil {
		response.NotFound(c, "分类不存在")
		return
	}
	response.Success(c, adminres.ToCategoryResource(cat))
}

// Store 创建分类
func (ctrl *CategoryController) Store(c *gin.Context) {
	var req adminreq.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}
	cat := &models.ProductCategory{
		Name: req.Name, ParentID: req.ParentID,
		SortOrder: req.SortOrder, IsEnabled: isEnabled,
	}
	if err := ctrl.service.Create(cat); err != nil {
		response.InternalError(c, "创建失败")
		return
	}
	response.Created(c, adminres.ToCategoryResource(cat))
}

// Update 更新分类
func (ctrl *CategoryController) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cat, _ := ctrl.service.AdminGetByID(id)
	if cat == nil {
		response.NotFound(c, "分类不存在")
		return
	}
	var req adminreq.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
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
		response.InternalError(c, "更新失败")
		return
	}
	response.Updated(c, adminres.ToCategoryResource(cat))
}

// Destroy 软删除
func (ctrl *CategoryController) Destroy(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := ctrl.service.Delete(id); err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Deleted(c)
}

// Restore 恢复已删除
func (ctrl *CategoryController) Restore(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := ctrl.service.Restore(id); err != nil {
		response.InternalError(c, "恢复失败")
		return
	}
	response.SuccessWithMsg(c, "恢复成功", nil)
}
