// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
)

type Querier interface {
	CreateCharacter(ctx context.Context, arg CreateCharacterParams) (Character, error)
	CreateChildCategory(ctx context.Context, arg CreateChildCategoryParams) (ChildCategory, error)
	CreateImage(ctx context.Context, arg CreateImageParams) (Image, error)
	CreateImageCharacterRelations(ctx context.Context, arg CreateImageCharacterRelationsParams) (ImageCharactersRelation, error)
	CreateImageParentCategoryRelations(ctx context.Context, arg CreateImageParentCategoryRelationsParams) (ImageParentCategoriesRelation, error)
	CreateOperator(ctx context.Context, arg CreateOperatorParams) (Operator, error)
	CreateParentCategory(ctx context.Context, arg CreateParentCategoryParams) (ParentCategory, error)
	DeleteCharacter(ctx context.Context, id int64) error
	DeleteChildCategory(ctx context.Context, id int64) error
	DeleteImage(ctx context.Context, id int64) error
	DeleteImageCharacterRelations(ctx context.Context, id int64) error
	DeleteImageParentCategoryRelations(ctx context.Context, id int64) error
	DeleteParentCategory(ctx context.Context, id int64) error
	GetCharacter(ctx context.Context, id int64) (Character, error)
	GetChildCategoriesByParentID(ctx context.Context, parentID int64) ([]ChildCategory, error)
	GetChildCategory(ctx context.Context, id int64) (ChildCategory, error)
	GetImage(ctx context.Context, id int64) (Image, error)
	GetOperator(ctx context.Context, id int64) (Operator, error)
	GetParentCategory(ctx context.Context, id int64) (ParentCategory, error)
	ListCharacters(ctx context.Context, arg ListCharactersParams) ([]Character, error)
	ListChildCategories(ctx context.Context, arg ListChildCategoriesParams) ([]ChildCategory, error)
	ListImage(ctx context.Context, arg ListImageParams) ([]Image, error)
	ListImageCharacterRelationsByImageID(ctx context.Context, imageID int64) ([]ImageCharactersRelation, error)
	ListImageCharacterRelationsByParentCategoryID(ctx context.Context, characterID int64) ([]ImageCharactersRelation, error)
	ListImageParentCategoryRelationsByImageID(ctx context.Context, imageID int64) ([]ImageParentCategoriesRelation, error)
	ListImageParentCategoryRelationsByParentCategoryID(ctx context.Context, parentCategoryID int64) ([]ImageParentCategoriesRelation, error)
	ListParentCategories(ctx context.Context, arg ListParentCategoriesParams) ([]ParentCategory, error)
	UpdateCharacter(ctx context.Context, arg UpdateCharacterParams) (Character, error)
	UpdateChildCategory(ctx context.Context, arg UpdateChildCategoryParams) (ChildCategory, error)
	UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error)
	UpdateImageCharacterRelations(ctx context.Context, arg UpdateImageCharacterRelationsParams) (ImageCharactersRelation, error)
	UpdateImageParentCategoryRelations(ctx context.Context, arg UpdateImageParentCategoryRelationsParams) (ImageParentCategoriesRelation, error)
	UpdateOperator(ctx context.Context, arg UpdateOperatorParams) (Operator, error)
	UpdateParentCategory(ctx context.Context, arg UpdateParentCategoryParams) (ParentCategory, error)
}

var _ Querier = (*Queries)(nil)
