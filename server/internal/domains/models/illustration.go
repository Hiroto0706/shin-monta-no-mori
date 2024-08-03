package model

import (
	db "shin-monta-no-mori/internal/db/sqlc"
)

type (
	Illustration struct {
		Image      db.Image
		Characters []*Character
		Categories []*Category
	}
)

func NewIllustration() *Illustration {
	return &Illustration{
		Image:      db.Image{},
		Characters: []*Character{},
		Categories: []*Category{
			{
				ParentCategory: db.ParentCategory{},
				ChildCategory:  []db.ChildCategory{},
			},
		},
	}
}
