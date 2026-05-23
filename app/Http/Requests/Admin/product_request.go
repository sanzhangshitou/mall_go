package adminreq

// CreateProductRequest 创建商品
type CreateProductRequest struct {
	CategoryID  int64      `json:"category_id" binding:"required"`
	Name        string     `json:"name" binding:"required,max=200"`
	Subtitle    string     `json:"subtitle"`
	Description string     `json:"description"`
	MainImage   string     `json:"main_image"`
	Images      []string   `json:"images"`
	Status      int        `json:"status"`
	Specs       []SpecItem `json:"specs"`
	Unit        string     `json:"unit"`
	SortOrder   int        `json:"sort_order"`
}

// UpdateProductRequest 更新商品
type UpdateProductRequest struct {
	CategoryID  *int64      `json:"category_id"`
	Name        string      `json:"name" binding:"max=200"`
	Subtitle    *string     `json:"subtitle"`
	Description *string     `json:"description"`
	MainImage   *string     `json:"main_image"`
	Images      *[]string   `json:"images"`
	Status      *int        `json:"status"`
	Specs       *[]SpecItem `json:"specs"`
	Unit        *string     `json:"unit"`
	SortOrder   *int        `json:"sort_order"`
}

// ProductListRequest 商品列表查询 (管理后台可查已删除)
type ProductListRequest struct {
	CategoryID     int64  `form:"category_id"`
	Status         *int   `form:"status"`
	Keyword        string `form:"keyword"`
	SortField      string `form:"sort_field"`
	SortOrder      string `form:"sort_order"`
	IncludeDeleted bool   `form:"include_deleted"`
}

// UpdateStatusRequest 状态变更
type UpdateStatusRequest struct {
	Status int `json:"status" binding:"required,min=0,max=3"`
}

// SpecItem 规格项
type SpecItem struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
