// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
)

type Querier interface {
	CreateCharacter(ctx context.Context, arg CreateCharacterParams) (Character, error)
	CreateImage(ctx context.Context, arg CreateImageParams) (Image, error)
	DeleteCharacter(ctx context.Context, id int64) error
	DeleteImage(ctx context.Context, id int64) error
	GetCharacter(ctx context.Context, id int64) (Character, error)
	GetImage(ctx context.Context, id int64) (Image, error)
	ListCharacters(ctx context.Context, arg ListCharactersParams) ([]Character, error)
	ListImage(ctx context.Context, arg ListImageParams) ([]Image, error)
	UpdateCharacter(ctx context.Context, arg UpdateCharacterParams) (Character, error)
	UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error)
}

var _ Querier = (*Queries)(nil)
