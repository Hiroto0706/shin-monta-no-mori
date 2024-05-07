package model

import (
	db "shin-monta-no-mori/server/internal/db/sqlc"
)

type (
	Illustration struct {
		Image     db.Image
		Character []db.Character
		Category  []*Category
	}

	Category struct {
		ParentCategory db.ParentCategory
		ChildCategory  []db.ChildCategory
	}
)

func NewIllustration() *Illustration {
	return &Illustration{
		Image:     db.Image{},
		Character: []db.Character{},
		Category:  []*Category{},
	}
}

func NewCategory() *Category {
	return &Category{
		ParentCategory: db.ParentCategory{},
		ChildCategory:  []db.ChildCategory{},
	}
}
