package model

import db "shin-monta-no-mori/server/internal/db/sqlc"

type (
	Operator struct {
		Operator db.Operator
	}
)

func NewOperator() *Operator {
	return &Operator{
		Operator: db.Operator{},
	}
}
