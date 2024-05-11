// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: child_categories.sql

package db

import (
	"context"
	"time"
)

const createChildCategory = `-- name: CreateChildCategory :one
INSERT INTO child_categories (name, parent_id)
VALUES ($1, $2)
RETURNING id, name, parent_id, updated_at, created_at
`

type CreateChildCategoryParams struct {
	Name     string `json:"name"`
	ParentID int64  `json:"parent_id"`
}

func (q *Queries) CreateChildCategory(ctx context.Context, arg CreateChildCategoryParams) (ChildCategory, error) {
	row := q.db.QueryRowContext(ctx, createChildCategory, arg.Name, arg.ParentID)
	var i ChildCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ParentID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteChildCategory = `-- name: DeleteChildCategory :exec
DELETE FROM child_categories
WHERE id = $1
`

func (q *Queries) DeleteChildCategory(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteChildCategory, id)
	return err
}

const getChildCategoriesByParentID = `-- name: GetChildCategoriesByParentID :many
SELECT id, name, parent_id, updated_at, created_at
FROM child_categories
WHERE parent_id = $1
`

func (q *Queries) GetChildCategoriesByParentID(ctx context.Context, parentID int64) ([]ChildCategory, error) {
	rows, err := q.db.QueryContext(ctx, getChildCategoriesByParentID, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ChildCategory{}
	for rows.Next() {
		var i ChildCategory
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
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

const getChildCategory = `-- name: GetChildCategory :one
SELECT id, name, parent_id, updated_at, created_at
FROM child_categories
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetChildCategory(ctx context.Context, id int64) (ChildCategory, error) {
	row := q.db.QueryRowContext(ctx, getChildCategory, id)
	var i ChildCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ParentID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const listChildCategories = `-- name: ListChildCategories :many
SELECT id, name, parent_id, updated_at, created_at
FROM child_categories
ORDER BY id DESC
LIMIT $1 OFFSET $2
`

type ListChildCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListChildCategories(ctx context.Context, arg ListChildCategoriesParams) ([]ChildCategory, error) {
	rows, err := q.db.QueryContext(ctx, listChildCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ChildCategory{}
	for rows.Next() {
		var i ChildCategory
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
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

const updateChildCategory = `-- name: UpdateChildCategory :one
UPDATE child_categories
SET name = $2,
  parent_id = $3,
  updated_at = $4
WHERE id = $1
RETURNING id, name, parent_id, updated_at, created_at
`

type UpdateChildCategoryParams struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  int64     `json:"parent_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateChildCategory(ctx context.Context, arg UpdateChildCategoryParams) (ChildCategory, error) {
	row := q.db.QueryRowContext(ctx, updateChildCategory,
		arg.ID,
		arg.Name,
		arg.ParentID,
		arg.UpdatedAt,
	)
	var i ChildCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ParentID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
