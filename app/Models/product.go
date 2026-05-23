package models

import (
	"database/sql"
	"encoding/json"
)

// Product 商品主表
type Product struct {
	ID          int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CategoryID  int64           `gorm:"column:category_id;not null" json:"category_id"`
	Name        string          `gorm:"column:name;type:varchar(200);not null" json:"name"`
	Subtitle    sql.NullString  `gorm:"column:subtitle;type:varchar(500)" json:"subtitle"`
	Description sql.NullString  `gorm:"column:description;type:text" json:"description"`
	MainImage   sql.NullString  `gorm:"column:main_image;type:varchar(500)" json:"main_image"`
	Images      json.RawMessage `gorm:"column:images;type:json" json:"images"`
	Status      ProductStatus   `gorm:"column:status;default:0" json:"status"`
	Specs       json.RawMessage `gorm:"column:specs;type:json" json:"specs"`
	Unit        sql.NullString  `gorm:"column:unit;type:varchar(20)" json:"unit"`
	SortOrder   int             `gorm:"column:sort_order;default:0" json:"sort_order"`
	SalesCount  int             `gorm:"column:sales_count;default:0" json:"sales_count"`
	IsDeleted   bool            `gorm:"column:is_deleted;default:0" json:"-"`
	CreatedAt   sql.NullTime    `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   sql.NullTime    `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// 关联 (手动填充, 非数据库字段)
	Skus      []ProductSku   `gorm:"-" json:"skus,omitempty"`
	ImageList []ProductImage `gorm:"-" json:"image_list,omitempty"`
}

func (Product) TableName() string {
	return "product"
}
