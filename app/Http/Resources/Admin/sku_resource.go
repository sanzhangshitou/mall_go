package adminres

import (
	"encoding/json"

	models "mall/app/Models"
)

// SkuResource SKU 资源（管理后台含 is_deleted）
type SkuResource struct {
	ID          int64           `json:"id"`
	ProductID   int64           `json:"product_id"`
	SkuCode     string          `json:"sku_code"`
	SpecValues  json.RawMessage `json:"spec_values"`
	Price       float64         `json:"price"`
	MarketPrice *float64        `json:"market_price"`
	Stock       int             `json:"stock"`
	Image       string          `json:"image"`
	IsEnabled   bool            `json:"is_enabled"`
	IsDeleted   bool            `json:"is_deleted"`
	SortOrder   int             `json:"sort_order"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

func ToSkuResource(sku *models.ProductSku) *SkuResource {
	if sku == nil {
		return nil
	}
	r := &SkuResource{
		ID:         sku.ID,
		ProductID:  sku.ProductID,
		SpecValues: sku.SpecValues,
		Price:      sku.Price,
		Stock:      sku.Stock,
		IsEnabled:  sku.IsEnabled,
		IsDeleted:  sku.IsDeleted,
		SortOrder:  sku.SortOrder,
	}
	if sku.SkuCode.Valid {
		r.SkuCode = sku.SkuCode.String
	}
	if sku.MarketPrice.Valid {
		r.MarketPrice = &sku.MarketPrice.Float64
	}
	if sku.Image.Valid {
		r.Image = sku.Image.String
	}
	if sku.CreatedAt.Valid {
		r.CreatedAt = sku.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if sku.UpdatedAt.Valid {
		r.UpdatedAt = sku.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return r
}

func ToSkuList(skus []models.ProductSku) []*SkuResource {
	list := make([]*SkuResource, 0, len(skus))
	for i := range skus {
		list = append(list, ToSkuResource(&skus[i]))
	}
	return list
}

// ---------- 图片 ----------

type ImageResource struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	SkuID     *int64 `json:"sku_id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

func ToImageResource(img *models.ProductImage) *ImageResource {
	if img == nil {
		return nil
	}
	r := &ImageResource{
		ID:        img.ID,
		ProductID: img.ProductID,
		URL:       img.URL,
		SortOrder: img.SortOrder,
	}
	if img.SkuID.Valid {
		r.SkuID = &img.SkuID.Int64
	}
	return r
}

func ToImageList(images []models.ProductImage) []*ImageResource {
	list := make([]*ImageResource, 0, len(images))
	for i := range images {
		list = append(list, ToImageResource(&images[i]))
	}
	return list
}
