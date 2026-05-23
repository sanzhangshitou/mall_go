package services

import (
	"errors"
	"fmt"

	models "mall/app/Models"
	"mall/app/Support/database"

	"gorm.io/gorm"
)

type ProductService struct {
	categoryService *CategoryService
	skuService      *ProductSkuService
	imageService    *ProductImageService
}

func NewProductService() *ProductService {
	return &ProductService{
		categoryService: NewCategoryService(),
		skuService:      NewProductSkuService(),
		imageService:    NewProductImageService(),
	}
}

// List 商品列表 (符合规范: 先查主表, 再批量补分类名)
func (s *ProductService) List(query ProductListQuery) ([]models.Product, int64, error) {
	db := database.DB().Model(&models.Product{}).Where("is_deleted = ?", false)

	// -- 筛选条件 --
	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}

	// 总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	// 排序
	orderClause := "sort_order ASC, id DESC"
	if query.SortField != "" {
		direction := "ASC"
		if query.SortOrder == "desc" {
			direction = "DESC"
		}
		allowed := map[string]bool{"sales_count": true, "created_at": true, "price": true}
		if allowed[query.SortField] {
			orderClause = fmt.Sprintf("%s %s", query.SortField, direction)
		}
	}

	// 分页查询主表
	var products []models.Product
	if err := db.Order(orderClause).
		Offset(query.Offset).
		Limit(query.PageSize).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetByID 获取商品详情 (符合规范: 先查主表, 再补 SKU+图片)
func (s *ProductService) GetByID(id int64) (*models.Product, error) {
	var product models.Product
	err := database.DB().
		Where("id = ?", id).
		Where("is_deleted = ?", false).
		First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// 批量补 SKU
	skus, _ := s.skuService.GetByProductID(product.ID)
	product.Skus = skus

	// 批量补图片
	images, _ := s.imageService.GetByProductID(product.ID)
	product.ImageList = images

	return &product, nil
}

// Create 创建商品 (仅主表, SKU/图片由调用方单独创建)
func (s *ProductService) Create(product *models.Product) error {
	return database.DB().Create(product).Error
}

// Update 更新商品
func (s *ProductService) Update(product *models.Product) error {
	return database.DB().Model(product).Omit("id", "created_at").Updates(product).Error
}

// Delete 软删除商品
func (s *ProductService) Delete(id int64) error {
	return database.DB().
		Model(&models.Product{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// ---------- Admin ----------

// AdminList 管理后台商品列表（可选包含已删除）
func (s *ProductService) AdminList(query ProductListQuery, includeDeleted bool) ([]models.Product, int64, error) {
	db := database.DB().Model(&models.Product{})
	if !includeDeleted {
		db = db.Where("is_deleted = ?", false)
	}

	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	orderClause := "sort_order ASC, id DESC"
	if query.SortField != "" {
		direction := "ASC"
		if query.SortOrder == "desc" {
			direction = "DESC"
		}
		if map[string]bool{"sales_count": true, "created_at": true}[query.SortField] {
			orderClause = fmt.Sprintf("%s %s", query.SortField, direction)
		}
	}

	var products []models.Product
	if err := db.Order(orderClause).Offset(query.Offset).Limit(query.PageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

// AdminGetByID 管理后台获取商品详情（含已删除）
func (s *ProductService) AdminGetByID(id int64) (*models.Product, error) {
	var product models.Product
	err := database.DB().Where("id = ?", id).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	skus, _ := s.skuService.AdminGetByProductID(product.ID)
	product.Skus = skus
	images, _ := s.imageService.GetByProductID(product.ID)
	product.ImageList = images
	return &product, nil
}

// Restore 恢复软删除
func (s *ProductService) Restore(id int64) error {
	return database.DB().Model(&models.Product{}).Where("id = ?", id).Update("is_deleted", false).Error
}

// UpdateStatus 修改商品状态
func (s *ProductService) UpdateStatus(id int64, status models.ProductStatus) error {
	return database.DB().Model(&models.Product{}).Where("id = ?", id).Update("status", status).Error
}

// ---------- 查询参数 ----------

type ProductListQuery struct {
	CategoryID int64
	Status     *int
	Keyword    string
	SortField  string
	SortOrder  string
	Page       int
	PageSize   int
	Offset     int
}
