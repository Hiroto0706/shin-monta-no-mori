// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	CountImages(ctx context.Context) (int64, error)
	CreateCharacter(ctx context.Context, arg CreateCharacterParams) (Character, error)
	CreateChildCategory(ctx context.Context, arg CreateChildCategoryParams) (ChildCategory, error)
	CreateImage(ctx context.Context, arg CreateImageParams) (Image, error)
	CreateImageCharacterRelations(ctx context.Context, arg CreateImageCharacterRelationsParams) (ImageCharactersRelation, error)
	CreateImageChildCategoryRelations(ctx context.Context, arg CreateImageChildCategoryRelationsParams) (ImageChildCategoriesRelation, error)
	CreateImageParentCategoryRelations(ctx context.Context, arg CreateImageParentCategoryRelationsParams) (ImageParentCategoriesRelation, error)
	CreateOperator(ctx context.Context, arg CreateOperatorParams) (Operator, error)
	CreateParentCategory(ctx context.Context, arg CreateParentCategoryParams) (ParentCategory, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	DeleteAllChildCategoriesByParentCategoryID(ctx context.Context, parentID int64) error
	DeleteAllImageCharacterRelationsByCharacterID(ctx context.Context, characterID int64) error
	DeleteAllImageCharacterRelationsByImageID(ctx context.Context, imageID int64) error
	DeleteAllImageChildCategoryRelationsByChildCategoryID(ctx context.Context, childCategoryID int64) error
	DeleteAllImageChildCategoryRelationsByImageID(ctx context.Context, imageID int64) error
	DeleteAllImageParentCategoryRelationsByImageID(ctx context.Context, imageID int64) error
	DeleteAllImageParentCategoryRelationsByParentCategoryID(ctx context.Context, parentCategoryID int64) error
	DeleteCharacter(ctx context.Context, id int64) error
	DeleteChildCategory(ctx context.Context, id int64) error
	DeleteImage(ctx context.Context, id int64) error
	DeleteImageCharacterRelations(ctx context.Context, id int64) error
	DeleteImageChildCategoryRelations(ctx context.Context, id int64) error
	DeleteImageParentCategoryRelations(ctx context.Context, id int64) error
	DeleteParentCategory(ctx context.Context, id int64) error
	FetchRandomImage(ctx context.Context, limit int32) ([]Image, error)
	GetCharacter(ctx context.Context, id int64) (Character, error)
	GetChildCategoriesByParentID(ctx context.Context, parentID int64) ([]ChildCategory, error)
	GetChildCategory(ctx context.Context, id int64) (ChildCategory, error)
	GetImage(ctx context.Context, id int64) (Image, error)
	GetOperatorByEmail(ctx context.Context, email string) (Operator, error)
	GetParentCategory(ctx context.Context, id int64) (ParentCategory, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	ListAllCharacters(ctx context.Context) ([]Character, error)
	ListCharacters(ctx context.Context, arg ListCharactersParams) ([]Character, error)
	ListChildCategories(ctx context.Context, arg ListChildCategoriesParams) ([]ChildCategory, error)
	ListImage(ctx context.Context, arg ListImageParams) ([]Image, error)
	ListImageCharacterRelationsByCharacterIDWIthPagination(ctx context.Context, arg ListImageCharacterRelationsByCharacterIDWIthPaginationParams) ([]ImageCharactersRelation, error)
	ListImageCharacterRelationsByImageID(ctx context.Context, imageID int64) ([]ImageCharactersRelation, error)
	ListImageChildCategoryRelationsByChildCategoryID(ctx context.Context, childCategoryID int64) ([]ImageChildCategoriesRelation, error)
	ListImageChildCategoryRelationsByChildCategoryIDWithPagination(ctx context.Context, arg ListImageChildCategoryRelationsByChildCategoryIDWithPaginationParams) ([]ImageChildCategoriesRelation, error)
	ListImageChildCategoryRelationsByImageID(ctx context.Context, imageID int64) ([]ImageChildCategoriesRelation, error)
	ListImageParentCategoryRelationsByImageID(ctx context.Context, imageID int64) ([]ImageParentCategoriesRelation, error)
	ListImageParentCategoryRelationsByParentCategoryID(ctx context.Context, parentCategoryID int64) ([]ImageParentCategoriesRelation, error)
	ListImageParentCategoryRelationsByParentCategoryIDWithPagination(ctx context.Context, arg ListImageParentCategoryRelationsByParentCategoryIDWithPaginationParams) ([]ImageParentCategoriesRelation, error)
	ListParentCategories(ctx context.Context) ([]ParentCategory, error)
	SearchCharacters(ctx context.Context, arg SearchCharactersParams) ([]Character, error)
	SearchImages(ctx context.Context, arg SearchImagesParams) ([]Image, error)
	SearchParentCategories(ctx context.Context, query sql.NullString) ([]ParentCategory, error)
	UpdateCharacter(ctx context.Context, arg UpdateCharacterParams) (Character, error)
	UpdateChildCategory(ctx context.Context, arg UpdateChildCategoryParams) (ChildCategory, error)
	UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error)
	UpdateImageCharacterRelations(ctx context.Context, arg UpdateImageCharacterRelationsParams) (ImageCharactersRelation, error)
	UpdateImageChildCategoryRelations(ctx context.Context, arg UpdateImageChildCategoryRelationsParams) (ImageChildCategoriesRelation, error)
	UpdateImageParentCategoryRelations(ctx context.Context, arg UpdateImageParentCategoryRelationsParams) (ImageParentCategoriesRelation, error)
	UpdateOperator(ctx context.Context, arg UpdateOperatorParams) (Operator, error)
	UpdateParentCategory(ctx context.Context, arg UpdateParentCategoryParams) (ParentCategory, error)
}

var _ Querier = (*Queries)(nil)
