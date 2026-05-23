package apires

import (
	"encoding/json"

	models "mall/app/Models"
)

// ProductResource 商品资源
type ProductResource struct {
	ID          int64           `json:"id"`
	CategoryID  int64           `json:"category_id"`
	Name        string          `json:"name"`
	Subtitle    string          `json:"subtitle"`
	Description string          `json:"description"`
	MainImage   string          `json:"main_image"`
	Images      json.RawMessage `json:"images"`
	Status      int             `json:"status"`
	StatusLabel string          `json:"status_label"`
	Specs       json.RawMessage `json:"specs"`
	Unit        string          `json:"unit"`
	SortOrder   int             `json:"sort_order"`
	SalesCount  int             `json:"sales_count"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`

	// 补充信息
	CategoryName string `json:"category_name,omitempty"`
}

// ToProductResource 单条转换
func ToProductResource(p *models.Product) *ProductResource {
	if p == nil {
		return nil
	}
	r := &ProductResource{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Images:      p.Images,
		Status:      int(p.Status),
		StatusLabel: models.ProductStatus(p.Status).Label(),
		Specs:       p.Specs,
		SortOrder:   p.SortOrder,
		SalesCount:  p.SalesCount,
	}
	if p.Subtitle.Valid {
		r.Subtitle = p.Subtitle.String
	}
	if p.Description.Valid {
		r.Description = p.Description.String
	}
	if p.MainImage.Valid {
		r.MainImage = p.MainImage.String
	}
	if p.Unit.Valid {
		r.Unit = p.Unit.String
	}
	if p.CreatedAt.Valid {
		r.CreatedAt = p.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if p.UpdatedAt.Valid {
		r.UpdatedAt = p.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return r
}

// ToProductList 批量转换
func ToProductList(products []models.Product) []*ProductResource {
	list := make([]*ProductResource, 0, len(products))
	for i := range products {
		list = append(list, ToProductResource(&products[i]))
	}
	return list
}

// ---------- 商品详情 (含 SKU + 图片) ----------

// ProductDetailResource 商品详情资源
type ProductDetailResource struct {
	*ProductResource
	Skus   []*SkuResource   `json:"skus"`
	Images []*ImageResource `json:"images"`
}

// ToProductDetail 转换为详情
func ToProductDetail(p *models.Product) *ProductDetailResource {
	if p == nil {
		return nil
	}
	detail := &ProductDetailResource{
		ProductResource: ToProductResource(p),
	}

	// SKU 列表
	if p.Skus != nil {
		detail.Skus = ToSkuList(p.Skus)
	}

	// 图片列表
	if p.ImageList != nil {
		detail.Images = ToImageList(p.ImageList)
	}

	return detail
}
