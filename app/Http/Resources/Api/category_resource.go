package apires

import (
	models "mall/app/Models"
)

// CategoryResource 分类资源
type CategoryResource struct {
	ID        int64               `json:"id"`
	Name      string              `json:"name"`
	ParentID  int64               `json:"parent_id"`
	SortOrder int                 `json:"sort_order"`
	IsEnabled bool                `json:"is_enabled"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
	Children  []*CategoryResource `json:"children,omitempty"`
}

// ToCategoryResource 单条转换
func ToCategoryResource(cat *models.ProductCategory) *CategoryResource {
	if cat == nil {
		return nil
	}
	r := &CategoryResource{
		ID:        cat.ID,
		Name:      cat.Name,
		ParentID:  cat.ParentID,
		SortOrder: cat.SortOrder,
		IsEnabled: cat.IsEnabled,
	}
	if cat.CreatedAt.Valid {
		r.CreatedAt = cat.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if cat.UpdatedAt.Valid {
		r.UpdatedAt = cat.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return r
}

// ToCategoryTree 转换为树形结构
func ToCategoryTree(categories []*models.ProductCategory) []*CategoryResource {
	// 建立 ID -> Resource 的映射
	resourceMap := make(map[int64]*CategoryResource)
	for _, cat := range categories {
		resourceMap[cat.ID] = ToCategoryResource(cat)
	}

	// 构建树
	var roots []*CategoryResource
	for _, cat := range categories {
		r := resourceMap[cat.ID]
		if cat.ParentID == 0 {
			roots = append(roots, r)
		} else if parent, ok := resourceMap[cat.ParentID]; ok {
			parent.Children = append(parent.Children, r)
		}
	}

	return roots
}
