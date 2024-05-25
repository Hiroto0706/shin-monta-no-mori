package model

import (
	db "shin-monta-no-mori/server/internal/db/sqlc"
)

type (
	Illustration struct {
		Image     db.Image
		Character []*Character
		Category  []*Category
	}
)

func NewIllustration() *Illustration {
	return &Illustration{
		Image:     db.Image{},
		Character: []*Character{},
		Category: []*Category{
			{
				ParentCategory: db.ParentCategory{},
				ChildCategory:  []db.ChildCategory{},
			},
		},
	}
}
