// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: image_parent_categories.sql

package db

import (
	"context"
)

const createImageParentCategoryRelations = `-- name: CreateImageParentCategoryRelations :one
INSERT INTO image_parent_categories_relations (image_id, parent_category_id)
VALUES ($1, $2)
RETURNING id, image_id, parent_category_id
`

type CreateImageParentCategoryRelationsParams struct {
	ImageID          int64 `json:"image_id"`
	ParentCategoryID int64 `json:"parent_category_id"`
}

func (q *Queries) CreateImageParentCategoryRelations(ctx context.Context, arg CreateImageParentCategoryRelationsParams) (ImageParentCategoriesRelation, error) {
	row := q.db.QueryRowContext(ctx, createImageParentCategoryRelations, arg.ImageID, arg.ParentCategoryID)
	var i ImageParentCategoriesRelation
	err := row.Scan(&i.ID, &i.ImageID, &i.ParentCategoryID)
	return i, err
}

const deleteAllImageParentCategoryRelationsByImageID = `-- name: DeleteAllImageParentCategoryRelationsByImageID :exec
DELETE FROM image_parent_categories_relations
WHERE image_id = $1
`

func (q *Queries) DeleteAllImageParentCategoryRelationsByImageID(ctx context.Context, imageID int64) error {
	_, err := q.db.ExecContext(ctx, deleteAllImageParentCategoryRelationsByImageID, imageID)
	return err
}

const deleteAllImageParentCategoryRelationsByParentCategoryID = `-- name: DeleteAllImageParentCategoryRelationsByParentCategoryID :exec
DELETE FROM image_parent_categories_relations
WHERE parent_category_id = $1
`

func (q *Queries) DeleteAllImageParentCategoryRelationsByParentCategoryID(ctx context.Context, parentCategoryID int64) error {
	_, err := q.db.ExecContext(ctx, deleteAllImageParentCategoryRelationsByParentCategoryID, parentCategoryID)
	return err
}

const deleteImageParentCategoryRelations = `-- name: DeleteImageParentCategoryRelations :exec
DELETE FROM image_parent_categories_relations
WHERE id = $1
`

func (q *Queries) DeleteImageParentCategoryRelations(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteImageParentCategoryRelations, id)
	return err
}

const listImageParentCategoryRelationsByImageID = `-- name: ListImageParentCategoryRelationsByImageID :many
SELECT id, image_id, parent_category_id
FROM image_parent_categories_relations
WHERE image_id = $1
ORDER BY image_id DESC
`

func (q *Queries) ListImageParentCategoryRelationsByImageID(ctx context.Context, imageID int64) ([]ImageParentCategoriesRelation, error) {
	rows, err := q.db.QueryContext(ctx, listImageParentCategoryRelationsByImageID, imageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ImageParentCategoriesRelation{}
	for rows.Next() {
		var i ImageParentCategoriesRelation
		if err := rows.Scan(&i.ID, &i.ImageID, &i.ParentCategoryID); err != nil {
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

const listImageParentCategoryRelationsByParentCategoryID = `-- name: ListImageParentCategoryRelationsByParentCategoryID :many
SELECT id, image_id, parent_category_id
FROM image_parent_categories_relations
WHERE parent_category_id = $1
ORDER BY parent_category_id DESC
`

func (q *Queries) ListImageParentCategoryRelationsByParentCategoryID(ctx context.Context, parentCategoryID int64) ([]ImageParentCategoriesRelation, error) {
	rows, err := q.db.QueryContext(ctx, listImageParentCategoryRelationsByParentCategoryID, parentCategoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ImageParentCategoriesRelation{}
	for rows.Next() {
		var i ImageParentCategoriesRelation
		if err := rows.Scan(&i.ID, &i.ImageID, &i.ParentCategoryID); err != nil {
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

const listImageParentCategoryRelationsByParentCategoryIDWithPagination = `-- name: ListImageParentCategoryRelationsByParentCategoryIDWithPagination :many
SELECT id, image_id, parent_category_id
FROM image_parent_categories_relations
WHERE parent_category_id = $3
ORDER BY parent_category_id DESC
LIMIT $1 OFFSET $2
`

type ListImageParentCategoryRelationsByParentCategoryIDWithPaginationParams struct {
	Limit            int32 `json:"limit"`
	Offset           int32 `json:"offset"`
	ParentCategoryID int64 `json:"parent_category_id"`
}

func (q *Queries) ListImageParentCategoryRelationsByParentCategoryIDWithPagination(ctx context.Context, arg ListImageParentCategoryRelationsByParentCategoryIDWithPaginationParams) ([]ImageParentCategoriesRelation, error) {
	rows, err := q.db.QueryContext(ctx, listImageParentCategoryRelationsByParentCategoryIDWithPagination, arg.Limit, arg.Offset, arg.ParentCategoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ImageParentCategoriesRelation{}
	for rows.Next() {
		var i ImageParentCategoriesRelation
		if err := rows.Scan(&i.ID, &i.ImageID, &i.ParentCategoryID); err != nil {
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

const updateImageParentCategoryRelations = `-- name: UpdateImageParentCategoryRelations :one
UPDATE image_parent_categories_relations
SET image_id = $2,
  parent_category_id = $3
WHERE id = $1
RETURNING id, image_id, parent_category_id
`

type UpdateImageParentCategoryRelationsParams struct {
	ID               int64 `json:"id"`
	ImageID          int64 `json:"image_id"`
	ParentCategoryID int64 `json:"parent_category_id"`
}

func (q *Queries) UpdateImageParentCategoryRelations(ctx context.Context, arg UpdateImageParentCategoryRelationsParams) (ImageParentCategoriesRelation, error) {
	row := q.db.QueryRowContext(ctx, updateImageParentCategoryRelations, arg.ID, arg.ImageID, arg.ParentCategoryID)
	var i ImageParentCategoriesRelation
	err := row.Scan(&i.ID, &i.ImageID, &i.ParentCategoryID)
	return i, err
}
