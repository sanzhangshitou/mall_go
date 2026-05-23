package adminreq

// CreateSkuRequest 创建SKU
type CreateSkuRequest struct {
	SkuCode     string            `json:"sku_code"`
	SpecValues  map[string]string `json:"spec_values"`
	Price       float64           `json:"price" binding:"required,gt=0"`
	MarketPrice *float64          `json:"market_price"`
	Stock       int               `json:"stock"`
	Image       string            `json:"image"`
	IsEnabled   *bool             `json:"is_enabled"`
	SortOrder   int               `json:"sort_order"`
}

// UpdateSkuRequest 更新SKU
type UpdateSkuRequest struct {
	SkuCode     *string            `json:"sku_code"`
	SpecValues  *map[string]string `json:"spec_values"`
	Price       *float64           `json:"price" binding:"omitempty,gt=0"`
	MarketPrice *float64           `json:"market_price"`
	Stock       *int               `json:"stock"`
	Image       *string            `json:"image"`
	IsEnabled   *bool              `json:"is_enabled"`
	SortOrder   *int               `json:"sort_order"`
}

// SkuListRequest SKU列表查询
type SkuListRequest struct {
	IncludeDeleted bool `form:"include_deleted"`
}
