package models

import "database/sql"

// ProductImage 商品图片表
type ProductImage struct {
	ID        int64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID int64         `gorm:"column:product_id;not null" json:"product_id"`
	SkuID     sql.NullInt64 `gorm:"column:sku_id" json:"sku_id"`
	URL       string        `gorm:"column:url;type:varchar(500);not null" json:"url"`
	SortOrder int           `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt sql.NullTime  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (ProductImage) TableName() string {
	return "product_image"
}
