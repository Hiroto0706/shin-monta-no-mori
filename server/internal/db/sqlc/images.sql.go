// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: images.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createImage = `-- name: CreateImage :one
INSERT INTO images (
    title,
    original_src,
    simple_src,
    original_filename,
    simple_filename
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, original_src, simple_src, updated_at, created_at, original_filename, simple_filename
`

type CreateImageParams struct {
	Title            string         `json:"title"`
	OriginalSrc      string         `json:"original_src"`
	SimpleSrc        sql.NullString `json:"simple_src"`
	OriginalFilename string         `json:"original_filename"`
	SimpleFilename   sql.NullString `json:"simple_filename"`
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, createImage,
		arg.Title,
		arg.OriginalSrc,
		arg.SimpleSrc,
		arg.OriginalFilename,
		arg.SimpleFilename,
	)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.OriginalSrc,
		&i.SimpleSrc,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.OriginalFilename,
		&i.SimpleFilename,
	)
	return i, err
}

const deleteImage = `-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1
`

func (q *Queries) DeleteImage(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteImage, id)
	return err
}

const getImage = `-- name: GetImage :one
SELECT id, title, original_src, simple_src, updated_at, created_at, original_filename, simple_filename
FROM images
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetImage(ctx context.Context, id int64) (Image, error) {
	row := q.db.QueryRowContext(ctx, getImage, id)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.OriginalSrc,
		&i.SimpleSrc,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.OriginalFilename,
		&i.SimpleFilename,
	)
	return i, err
}

const listImage = `-- name: ListImage :many
SELECT id, title, original_src, simple_src, updated_at, created_at, original_filename, simple_filename
FROM images
ORDER BY id DESC
LIMIT $1 OFFSET $2
`

type ListImageParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListImage(ctx context.Context, arg ListImageParams) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, listImage, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.OriginalSrc,
			&i.SimpleSrc,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.OriginalFilename,
			&i.SimpleFilename,
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

const searchImages = `-- name: SearchImages :many
SELECT DISTINCT id, title, original_src, simple_src, updated_at, created_at, original_filename, simple_filename
FROM images
WHERE title LIKE '%' || COALESCE($3) || '%'
  OR original_filename LIKE '%' || COALESCE($3) || '%'
ORDER BY id DESC
LIMIT $1 OFFSET $2
`

type SearchImagesParams struct {
	Limit  int32          `json:"limit"`
	Offset int32          `json:"offset"`
	Query  sql.NullString `json:"query"`
}

func (q *Queries) SearchImages(ctx context.Context, arg SearchImagesParams) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, searchImages, arg.Limit, arg.Offset, arg.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.OriginalSrc,
			&i.SimpleSrc,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.OriginalFilename,
			&i.SimpleFilename,
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

const updateImage = `-- name: UpdateImage :one
UPDATE images
SET title = $2,
  original_src = $3,
  simple_src = $4,
  original_filename = $5,
  simple_filename = $6,
  updated_at = $7
WHERE id = $1
RETURNING id, title, original_src, simple_src, updated_at, created_at, original_filename, simple_filename
`

type UpdateImageParams struct {
	ID               int64          `json:"id"`
	Title            string         `json:"title"`
	OriginalSrc      string         `json:"original_src"`
	SimpleSrc        sql.NullString `json:"simple_src"`
	OriginalFilename string         `json:"original_filename"`
	SimpleFilename   sql.NullString `json:"simple_filename"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (q *Queries) UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, updateImage,
		arg.ID,
		arg.Title,
		arg.OriginalSrc,
		arg.SimpleSrc,
		arg.OriginalFilename,
		arg.SimpleFilename,
		arg.UpdatedAt,
	)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.OriginalSrc,
		&i.SimpleSrc,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.OriginalFilename,
		&i.SimpleFilename,
	)
	return i, err
}
