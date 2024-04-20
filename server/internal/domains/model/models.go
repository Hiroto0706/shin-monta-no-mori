package model

import (
	db "shin-monta-no-mori/server/internal/db/sqlc"
)

type (
	Illustration struct {
		Image          db.Image
		Character      db.Character
		ParentCategory db.ParentCategory
		ChildCategory  db.ChildCategory
	}
)
