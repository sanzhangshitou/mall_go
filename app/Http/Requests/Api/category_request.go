package apireq

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	ParentID  int64  `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
	IsEnabled *bool  `json:"is_enabled"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name      string `json:"name" binding:"max=100"`
	ParentID  *int64 `json:"parent_id"`
	SortOrder *int   `json:"sort_order"`
	IsEnabled *bool  `json:"is_enabled"`
}
