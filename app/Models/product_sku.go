package models

import (
	"database/sql"
	"encoding/json"
)

// ProductSku 商品SKU表
type ProductSku struct {
	ID          int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID   int64           `gorm:"column:product_id;not null" json:"product_id"`
	SkuCode     sql.NullString  `gorm:"column:sku_code;type:varchar(100)" json:"sku_code"`
	SpecValues  json.RawMessage `gorm:"column:spec_values;type:json" json:"spec_values"`
	Price       float64         `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	MarketPrice sql.NullFloat64 `gorm:"column:market_price;type:decimal(10,2)" json:"market_price"`
	Stock       int             `gorm:"column:stock;default:0" json:"stock"`
	Image       sql.NullString  `gorm:"column:image;type:varchar(500)" json:"image"`
	IsEnabled   bool            `gorm:"column:is_enabled;default:1" json:"is_enabled"`
	SortOrder   int             `gorm:"column:sort_order;default:0" json:"sort_order"`
	IsDeleted   bool            `gorm:"column:is_deleted;default:0" json:"-"`
	CreatedAt   sql.NullTime    `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   sql.NullTime    `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ProductSku) TableName() string {
	return "product_sku"
}
