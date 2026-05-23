package models

import "database/sql"

// ProductCategory 商品分类表
type ProductCategory struct {
	ID        int64        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string       `gorm:"column:name;type:varchar(100);not null" json:"name"`
	ParentID  int64        `gorm:"column:parent_id;default:0" json:"parent_id"`
	SortOrder int          `gorm:"column:sort_order;default:0" json:"sort_order"`
	IsEnabled bool         `gorm:"column:is_enabled;default:1" json:"is_enabled"`
	IsDeleted bool         `gorm:"column:is_deleted;default:0" json:"-"`
	CreatedAt sql.NullTime `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt sql.NullTime `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// 非数据库字段 - 子分类树
	Children []*ProductCategory `gorm:"-" json:"children,omitempty"`
}

func (ProductCategory) TableName() string {
	return "product_category"
}
