package services

import (
	"errors"

	models "mall/app/Models"
	"mall/app/Support/database"

	"gorm.io/gorm"
)

type ProductImageService struct{}

func NewProductImageService() *ProductImageService {
	return &ProductImageService{}
}

// GetByProductID 获取商品下所有图片
func (s *ProductImageService) GetByProductID(productID int64) ([]models.ProductImage, error) {
	var images []models.ProductImage
	err := database.DB().
		Where("product_id = ?", productID).
		Order("sort_order ASC").
		Find(&images).Error
	return images, err
}

// GetByID 根据 ID 获取图片
func (s *ProductImageService) GetByID(id int64) (*models.ProductImage, error) {
	var img models.ProductImage
	err := database.DB().Where("id = ?", id).First(&img).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &img, err
}

// Create 创建图片
func (s *ProductImageService) Create(img *models.ProductImage) error {
	return database.DB().Create(img).Error
}

// Delete 删除图片
func (s *ProductImageService) Delete(id int64) error {
	return database.DB().Where("id = ?", id).Delete(&models.ProductImage{}).Error
}

// BatchCreate 批量创建
func (s *ProductImageService) BatchCreate(images []*models.ProductImage) error {
	return database.DB().Create(&images).Error
}
