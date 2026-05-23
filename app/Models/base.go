package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型 (不含主键, 由子表嵌入)
type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// SoftDeleteModel 带软删除的基础模型
type SoftDeleteModel struct {
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:is_deleted;index" json:"-"`
}
