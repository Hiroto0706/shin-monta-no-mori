package user

import (
	"fmt"
	"net/http"
	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
)

const (
	IMAGE_TYPE_CATEGORY = "category"
)

// TODO: 将来的にpager機能を持たせた方がいいかも？
type listCategoriesRequest struct {
	Page int64 `form:"p"`
}

type listCategoriesResponse struct {
	Categories []model.Category `json:"categories"`
}

// ListCategories godoc
// @Summary List categories
// @Description Retrieves a list of parent categories along with their child categories.
// @Accept  json
// @Produce  json
// @Success 200 {array} model/Category "A list of categories with parent and child category details."
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: The request is malformed or missing required fields."
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: Child categories not found for one or more parent categories."
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// @Router /api/v1/categories/list [get]
func ListCategories(ctx *app.AppContext) {
	// // TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listCategoriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}

	pcates, err := ctx.Server.Store.ListParentCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ctx.Server.Store.ListParentCategories : %w", err)))
		return
	}

	categories := make([]model.Category, len(pcates))
	for i, pcate := range pcates {
		ccates, err := ctx.Server.Store.GetChildCategoriesByParentID(ctx, pcate.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID : %w", err)))
			return
		}

		categories[i] = model.Category{
			ParentCategory: pcate,
			ChildCategory:  ccates,
		}
	}

	ctx.JSON(http.StatusOK, listCategoriesResponse{Categories: categories})
}

type listChildCategoriesResponse struct {
	ChildCategories []db.ChildCategory `json:"child_categories"`
}

// ListChildCategories は子カテゴリのリストを取得してクライアントに返します。
// @Summary 子カテゴリのリストを取得
// @Description 子カテゴリのリストを取得し、JSON形式でクライアントに返します。
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} listChildCategoriesResponse
// @Failure 500 {object} app.ErrorResponse "内部サーバーエラー"
// @Router /api/v1/categories/child/list [get]
func ListChildCategories(ctx *app.AppContext) {
	const FetchLimit = 5
	arg := db.ListChildCategoriesParams{
		Limit:  FetchLimit,
		Offset: 0,
	}
	childCategories, err := ctx.Server.Store.ListChildCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListChildCategories : %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, listChildCategoriesResponse{ChildCategories: childCategories})
}
