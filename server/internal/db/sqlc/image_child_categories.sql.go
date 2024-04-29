// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: image_child_categories.sql

package db

import (
	"context"
)

const createImageChildCategoryRelations = `-- name: CreateImageChildCategoryRelations :one
INSERT INTO image_child_categories_relations (image_id, child_category_id)
VALUES ($1, $2)
RETURNING id, image_id, child_category_id
`

type CreateImageChildCategoryRelationsParams struct {
	ImageID         int64 `json:"image_id"`
	ChildCategoryID int64 `json:"child_category_id"`
}

func (q *Queries) CreateImageChildCategoryRelations(ctx context.Context, arg CreateImageChildCategoryRelationsParams) (ImageChildCategoriesRelation, error) {
	row := q.db.QueryRowContext(ctx, createImageChildCategoryRelations, arg.ImageID, arg.ChildCategoryID)
	var i ImageChildCategoriesRelation
	err := row.Scan(&i.ID, &i.ImageID, &i.ChildCategoryID)
	return i, err
}

const deleteImageChildCategoryRelations = `-- name: DeleteImageChildCategoryRelations :exec
DELETE FROM image_child_categories_relations
WHERE id = $1
`

func (q *Queries) DeleteImageChildCategoryRelations(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteImageChildCategoryRelations, id)
	return err
}

const listImageChildCategoryRelationsByImageID = `-- name: ListImageChildCategoryRelationsByImageID :many
SELECT id, image_id, child_category_id
FROM image_child_categories_relations
WHERE image_id = $1
ORDER BY image_id DESC
`

func (q *Queries) ListImageChildCategoryRelationsByImageID(ctx context.Context, imageID int64) ([]ImageChildCategoriesRelation, error) {
	rows, err := q.db.QueryContext(ctx, listImageChildCategoryRelationsByImageID, imageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ImageChildCategoriesRelation{}
	for rows.Next() {
		var i ImageChildCategoriesRelation
		if err := rows.Scan(&i.ID, &i.ImageID, &i.ChildCategoryID); err != nil {
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

const listImageChildCategoryRelationsByParentCategoryID = `-- name: ListImageChildCategoryRelationsByParentCategoryID :many
SELECT id, image_id, child_category_id
FROM image_child_categories_relations
WHERE child_category_id = $1
ORDER BY child_category_id DESC
`

func (q *Queries) ListImageChildCategoryRelationsByParentCategoryID(ctx context.Context, childCategoryID int64) ([]ImageChildCategoriesRelation, error) {
	rows, err := q.db.QueryContext(ctx, listImageChildCategoryRelationsByParentCategoryID, childCategoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ImageChildCategoriesRelation{}
	for rows.Next() {
		var i ImageChildCategoriesRelation
		if err := rows.Scan(&i.ID, &i.ImageID, &i.ChildCategoryID); err != nil {
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

const updateImageChildCategoryRelations = `-- name: UpdateImageChildCategoryRelations :one
UPDATE image_child_categories_relations
SET image_id = $2,
  child_category_id = $3
WHERE id = $1
RETURNING id, image_id, child_category_id
`

type UpdateImageChildCategoryRelationsParams struct {
	ID              int64 `json:"id"`
	ImageID         int64 `json:"image_id"`
	ChildCategoryID int64 `json:"child_category_id"`
}

func (q *Queries) UpdateImageChildCategoryRelations(ctx context.Context, arg UpdateImageChildCategoryRelationsParams) (ImageChildCategoriesRelation, error) {
	row := q.db.QueryRowContext(ctx, updateImageChildCategoryRelations, arg.ID, arg.ImageID, arg.ChildCategoryID)
	var i ImageChildCategoriesRelation
	err := row.Scan(&i.ID, &i.ImageID, &i.ChildCategoryID)
	return i, err
}