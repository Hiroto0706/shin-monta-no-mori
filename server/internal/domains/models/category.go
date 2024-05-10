package model

import db "shin-monta-no-mori/server/internal/db/sqlc"

type (
	Category struct {
		ParentCategory db.ParentCategory
		ChildCategory  []db.ChildCategory
	}
)

func NewCategory() *Category {
	return &Category{
		ParentCategory: db.ParentCategory{},
		ChildCategory:  []db.ChildCategory{},
	}
}
