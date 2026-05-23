package services

import (
	"errors"

	models "mall/app/Models"
	"mall/app/Support/database"

	"gorm.io/gorm"
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// GetAll 获取全部分类 (平铺, 按 sort_order 排序)
func (s *CategoryService) GetAll() ([]*models.ProductCategory, error) {
	var categories []*models.ProductCategory
	err := database.DB().
		Where("is_enabled = ?", true).
		Where("is_deleted = ?", false).
		Order("sort_order ASC").
		Find(&categories).Error
	return categories, err
}

// GetByID 根据 ID 获取分类
func (s *CategoryService) GetByID(id int64) (*models.ProductCategory, error) {
	var cat models.ProductCategory
	err := database.DB().
		Where("id = ?", id).
		Where("is_deleted = ?", false).
		First(&cat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cat, err
}

// Create 创建分类
func (s *CategoryService) Create(cat *models.ProductCategory) error {
	return database.DB().Create(cat).Error
}

// Update 更新分类
func (s *CategoryService) Update(cat *models.ProductCategory) error {
	return database.DB().Model(cat).Omit("id", "created_at").Updates(cat).Error
}

// Delete 软删除分类
func (s *CategoryService) Delete(id int64) error {
	return database.DB().
		Model(&models.ProductCategory{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// GetByIDs 批量获取分类 (用于补数据)
func (s *CategoryService) GetByIDs(ids []int64) ([]models.ProductCategory, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var categories []models.ProductCategory
	err := database.DB().
		Where("id IN ?", ids).
		Where("is_deleted = ?", false).
		Find(&categories).Error
	return categories, err
}

// ---------- Admin ----------

// AdminGetAll 管理后台获取全部分类（可选包含禁用/已删除）
func (s *CategoryService) AdminGetAll(includeDisabled, includeDeleted bool) ([]*models.ProductCategory, error) {
	q := database.DB().Order("sort_order ASC")
	if !includeDisabled {
		q = q.Where("is_enabled = ?", true)
	}
	if !includeDeleted {
		q = q.Where("is_deleted = ?", false)
	}
	var categories []*models.ProductCategory
	err := q.Find(&categories).Error
	return categories, err
}

// AdminGetByID 管理后台获取分类（含已删除）
func (s *CategoryService) AdminGetByID(id int64) (*models.ProductCategory, error) {
	var cat models.ProductCategory
	err := database.DB().Where("id = ?", id).First(&cat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cat, err
}

// Restore 恢复软删除
func (s *CategoryService) Restore(id int64) error {
	return database.DB().
		Model(&models.ProductCategory{}).
		Where("id = ?", id).
		Update("is_deleted", false).Error
}
