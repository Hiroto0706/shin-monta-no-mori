// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: parent_categories.sql

package db

import (
	"context"
)

const createParentCategories = `-- name: CreateParentCategories :one
INSERT INTO parent_categories (name, src)
VALUES ($1, $2)
RETURNING id, name, src, updated_at, created_at
`

type CreateParentCategoriesParams struct {
	Name string `json:"name"`
	Src  string `json:"src"`
}

func (q *Queries) CreateParentCategories(ctx context.Context, arg CreateParentCategoriesParams) (ParentCategory, error) {
	row := q.db.QueryRowContext(ctx, createParentCategories, arg.Name, arg.Src)
	var i ParentCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Src,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteParentCategories = `-- name: DeleteParentCategories :exec
DELETE FROM parent_categories
WHERE id = $1
`

func (q *Queries) DeleteParentCategories(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteParentCategories, id)
	return err
}

const getParentCategories = `-- name: GetParentCategories :one
SELECT id, name, src, updated_at, created_at
FROM parent_categories
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetParentCategories(ctx context.Context, id int64) (ParentCategory, error) {
	row := q.db.QueryRowContext(ctx, getParentCategories, id)
	var i ParentCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Src,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const listParentCategories = `-- name: ListParentCategories :many
SELECT id, name, src, updated_at, created_at
FROM parent_categories
ORDER BY id DESC
LIMIT $1 OFFSET $2
`

type ListParentCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListParentCategories(ctx context.Context, arg ListParentCategoriesParams) ([]ParentCategory, error) {
	rows, err := q.db.QueryContext(ctx, listParentCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ParentCategory{}
	for rows.Next() {
		var i ParentCategory
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Src,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateParentCategories = `-- name: UpdateParentCategories :one
UPDATE parent_categories
SET name = $2,
  src = $3
WHERE id = $1
RETURNING id, name, src, updated_at, created_at
`

type UpdateParentCategoriesParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Src  string `json:"src"`
}

func (q *Queries) UpdateParentCategories(ctx context.Context, arg UpdateParentCategoriesParams) (ParentCategory, error) {
	row := q.db.QueryRowContext(ctx, updateParentCategories, arg.ID, arg.Name, arg.Src)
	var i ParentCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Src,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
