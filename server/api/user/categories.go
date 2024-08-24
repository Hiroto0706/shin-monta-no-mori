package user

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"shin-monta-no-mori/internal/app"
	"shin-monta-no-mori/internal/cache"
	db "shin-monta-no-mori/internal/db/sqlc"
	model "shin-monta-no-mori/internal/domains/models"
	"shin-monta-no-mori/pkg/lib/binder"
	"strconv"

	"github.com/redis/go-redis/v9"
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
	var req listCategoriesRequest
	if err := binder.BindQuery(ctx.Context, &req); err != nil {
		return
	}

	arg := db.ListParentCategoriesParams{
		Limit:  int32(ctx.Server.Config.CategoryFetchLimit),
		Offset: int32(int(req.Page) * ctx.Server.Config.CategoryFetchLimit),
	}
	pcates, err := ctx.Server.Store.ListParentCategories(ctx, arg)
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

// ListCategoriesAll handles the request to list all categories including their parent and child categories.
// @Summary List all categories
// @Description Get a list of all categories, including their parent and child categories.
// @Tags categories
// @Produce json
// @Success 200 {object} listCategoriesResponse
// @Failure 500 {object} app.ErrorResponse
// @param ctx AppContext
// @Router /categories/all [get]
func ListAllCategories(ctx *app.AppContext) {
	// TODO: redis周りの処理は関数化したい
	// Redisのキャッシュキーを設定
	cacheKey := cache.GetCategoriesAllKey()

	// Redisからキャッシュを取得
	var cachedResponse listCategoriesResponse
	err := ctx.Server.RedisClient.Get(ctx.Context, cacheKey, &cachedResponse)
	if err != nil && !errors.Is(err, redis.Nil) {
		// キャッシュの取得に失敗したが、デフォルトの動作としてDBからデータを取得する処理を続ける
		// TODO: loggerを追加する
		log.Println("failed to redis err : %w", err)
	}

	if err == nil {
		// キャッシュが存在する場合、それをレスポンスとして返す
		ctx.JSON(http.StatusOK, cachedResponse)
		return
	}

	pcates, err := ctx.Server.Store.ListAllParentCategories(ctx)
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

	if len(categories) > 0 {
		// レスポンスをキャッシュに保存
		response := listCategoriesResponse{
			Categories: categories,
		}
		// Redisへのセットが失敗しても処理を続行
		// TODO: loggerを追加する
		err = ctx.Server.RedisClient.Set(ctx.Context, cacheKey, response, cache.CacheDurationDay)
		if err != nil {
			log.Println("failed redis data set : %w", err)
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

type getChildCategoryResponse struct {
	ChildCategory db.ChildCategory `json:"child_category"`
}

// GetChildCategory は指定されたIDに基づいて子カテゴリを取得し、クライアントに返します。
//
// @Summary 特定の子カテゴリを取得
// @Description 指定されたIDに基づいて子カテゴリを取得し、JSON形式でクライアントに返します。
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "子カテゴリID"
// @Success 200 {object} getChildCategoryResponse
// @Failure 400 {object} app.ErrorResponse "無効なリクエストパラメータ"
// @Failure 500 {object} app.ErrorResponse "内部サーバーエラー"
// @Router /api/v1/categories/child/{id} [get]
func GetChildCategory(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}
	childCategory, err := ctx.Server.Store.GetChildCategory(ctx, int64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListChildCategories : %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, getChildCategoryResponse{ChildCategory: childCategory})
}
