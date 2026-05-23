package models

// ProductStatus 商品状态枚举
type ProductStatus int

const (
	ProductStatusDraft    ProductStatus = 0 // 草稿
	ProductStatusOnSale   ProductStatus = 1 // 上架
	ProductStatusOffSale  ProductStatus = 2 // 下架
	ProductStatusDisabled ProductStatus = 3 // 禁用
)

// StatusLabel 返回状态中文标签
func (s ProductStatus) Label() string {
	switch s {
	case ProductStatusDraft:
		return "草稿"
	case ProductStatusOnSale:
		return "上架"
	case ProductStatusOffSale:
		return "下架"
	case ProductStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}
