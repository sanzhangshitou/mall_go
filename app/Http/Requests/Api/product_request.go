package apireq

// CreateProductRequest 创建商品请求
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

// UpdateProductRequest 更新商品请求
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

// ProductListRequest 商品列表查询请求
type ProductListRequest struct {
	CategoryID int64  `form:"category_id"`
	Status     *int   `form:"status"`
	Keyword    string `form:"keyword"`
	SortField  string `form:"sort_field"`
	SortOrder  string `form:"sort_order"` // asc / desc
}

// SpecItem 规格项
type SpecItem struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
