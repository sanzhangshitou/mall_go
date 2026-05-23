package services

import (
	"encoding/json"
	"errors"
	"time"

	models "mall/app/Models"
	"mall/app/Support/cache"
	"mall/app/Support/database"

	"gorm.io/gorm"
)

const (
	cacheKeyCategories      = "categories:tree"
	cacheKeyAdminCategories = "categories:admin_tree"
	cacheTTL                = 10 * time.Minute
	cacheTTLAdmin           = 5 * time.Minute
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// GetAll 获取全部分类 (优先读缓存)
func (s *CategoryService) GetAll() ([]*models.ProductCategory, error) {
	key := cache.CacheKey(cacheKeyCategories)

	// 命中缓存
	if cached, err := cache.Get(key); err == nil && cached != "" {
		var categories []*models.ProductCategory
		if json.Unmarshal([]byte(cached), &categories) == nil {
			return categories, nil
		}
	}

	var categories []*models.ProductCategory
	err := database.DB().
		Where("is_enabled = ?", true).
		Where("is_deleted = ?", false).
		Order("sort_order ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(categories); err == nil {
		_ = cache.Set(key, data, cacheTTL)
	}

	return categories, nil
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
	if err := database.DB().Create(cat).Error; err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

// Update 更新分类
func (s *CategoryService) Update(cat *models.ProductCategory) error {
	if err := database.DB().Model(cat).Omit("id", "created_at").Updates(cat).Error; err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

// Delete 软删除分类
func (s *CategoryService) Delete(id int64) error {
	if err := database.DB().
		Model(&models.ProductCategory{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	s.invalidateCache()
	return nil
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

// AdminGetAll 管理后台获取全部分类（缓存全集，按需过滤）
func (s *CategoryService) AdminGetAll(includeDisabled, includeDeleted bool) ([]*models.ProductCategory, error) {
	key := cache.CacheKey(cacheKeyAdminCategories)

	var all []*models.ProductCategory

	if cached, err := cache.Get(key); err == nil && cached != "" {
		if json.Unmarshal([]byte(cached), &all) != nil {
			all = nil // 解析失败，走 DB
		}
	}

	if all == nil {
		// 缓存未命中，查全量
		if err := database.DB().Order("sort_order ASC").Find(&all).Error; err != nil {
			return nil, err
		}
		if data, err := json.Marshal(all); err == nil {
			_ = cache.Set(key, data, cacheTTLAdmin)
		}
	}

	// 内存过滤
	if includeDisabled && includeDeleted {
		return all, nil
	}
	result := make([]*models.ProductCategory, 0, len(all))
	for _, cat := range all {
		if !includeDisabled && !cat.IsEnabled {
			continue
		}
		if !includeDeleted && cat.IsDeleted {
			continue
		}
		result = append(result, cat)
	}
	return result, nil
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
	if err := database.DB().
		Model(&models.ProductCategory{}).
		Where("id = ?", id).
		Update("is_deleted", false).Error; err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

// invalidateCache 清除分类相关缓存
func (s *CategoryService) invalidateCache() {
	_ = cache.Del(
		cache.CacheKey(cacheKeyCategories),
		cache.CacheKey(cacheKeyAdminCategories),
	)
}
