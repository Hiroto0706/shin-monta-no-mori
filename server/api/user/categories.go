package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"shin-monta-no-mori/server/internal/app"
	model "shin-monta-no-mori/server/internal/domains/models"
)

const (
	IMAGE_TYPE_CATEGORY = "category"
)

// TODO: 将来的にpager機能を持たせた方がいいかも？
type listCategoriesRequest struct {
	Page int64 `form:"p"`
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
	// var req listCategoriesRequest
	// if err := c.ShouldBindQuery(&req); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
	// 	return
	// }

	pcates, err := ctx.Server.Store.ListParentCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ctx.Server.Store.ListParentCategories : %w", err)))
		return
	}

	categories := make([]model.Category, len(pcates))
	for i, pcate := range pcates {
		ccates, err := ctx.Server.Store.GetChildCategoriesByParentID(ctx, pcate.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID: %w", err)))
				return
			}

			ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID : %w", err)))
			return
		}

		categories[i] = model.Category{
			ParentCategory: pcate,
			ChildCategory:  ccates,
		}
	}

	ctx.JSON(http.StatusOK, categories)
}
