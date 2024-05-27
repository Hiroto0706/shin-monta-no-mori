package model

import db "shin-monta-no-mori/server/internal/db/sqlc"

type (
	Character struct {
		Character db.Character
	}
)

func NewCharacter() *Character {
	return &Character{
		Character: db.Character{},
	}
}
