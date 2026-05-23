package services

import (
	"errors"

	models "mall/app/Models"
	"mall/app/Support/database"

	"gorm.io/gorm"
)

type ProductSkuService struct{}

func NewProductSkuService() *ProductSkuService {
	return &ProductSkuService{}
}

// GetByProductID 获取商品下所有 SKU (含已删除)
func (s *ProductSkuService) GetByProductID(productID int64) ([]models.ProductSku, error) {
	var skus []models.ProductSku
	err := database.DB().
		Where("product_id = ?", productID).
		Where("is_deleted = ?", false).
		Order("sort_order ASC").
		Find(&skus).Error
	return skus, err
}

// GetByID 根据 ID 获取单个 SKU
func (s *ProductSkuService) GetByID(id int64) (*models.ProductSku, error) {
	var sku models.ProductSku
	err := database.DB().
		Where("id = ?", id).
		Where("is_deleted = ?", false).
		First(&sku).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &sku, err
}

// Create 创建 SKU
func (s *ProductSkuService) Create(sku *models.ProductSku) error {
	return database.DB().Create(sku).Error
}

// Update 更新 SKU
func (s *ProductSkuService) Update(sku *models.ProductSku) error {
	return database.DB().Model(sku).Omit("id", "product_id", "created_at").Updates(sku).Error
}

// Delete 软删除 SKU
func (s *ProductSkuService) Delete(id int64) error {
	return database.DB().
		Model(&models.ProductSku{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// BatchCreate 批量创建 SKU
func (s *ProductSkuService) BatchCreate(skus []*models.ProductSku) error {
	return database.DB().Create(&skus).Error
}

// ---------- Admin ----------

// AdminGetByProductID 管理后台获取商品 SKU（含已删除）
func (s *ProductSkuService) AdminGetByProductID(productID int64) ([]models.ProductSku, error) {
	var skus []models.ProductSku
	err := database.DB().
		Where("product_id = ?", productID).
		Order("sort_order ASC").
		Find(&skus).Error
	return skus, err
}

// AdminGetByID 管理后台获取 SKU（含已删除）
func (s *ProductSkuService) AdminGetByID(id int64) (*models.ProductSku, error) {
	var sku models.ProductSku
	err := database.DB().Where("id = ?", id).First(&sku).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &sku, err
}

// Restore 恢复软删除
func (s *ProductSkuService) Restore(id int64) error {
	return database.DB().
		Model(&models.ProductSku{}).
		Where("id = ?", id).
		Update("is_deleted", false).Error
}
