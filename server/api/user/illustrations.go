package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"shin-monta-no-mori/internal/app"
	"shin-monta-no-mori/internal/cache"
	db "shin-monta-no-mori/internal/db/sqlc"
	model "shin-monta-no-mori/internal/domains/models"
	"shin-monta-no-mori/internal/domains/service"
	"shin-monta-no-mori/pkg/lib/binder"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	IMAGE_TYPE_IMAGE = "image"
)

type listIllustrationsRequest struct {
	Page int64 `form:"p"`
}

type listIllustrationsResponse struct {
	Illustrations []*model.Illustration `json:"illustrations"`
}

// ListIllustrations godoc
// @Summary List illustrations
// @Description Retrieves a paginated list of illustrations based on the provided page number.
// @Accept  json
// @Produce  json
// @Param   p   query   int  true  "Page number for pagination"
// @Success 200 {array} model/Illustration "A list of illustrations"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: The request is malformed or missing required fields."
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// @Router /api/v1/illustrations/list [get]
func ListIllustrations(ctx *app.AppContext) {
	var req listIllustrationsRequest
	if err := binder.BindQuery(ctx.Context, &req); err != nil {
		return
	}

	// TODO: redis周りの処理は関数化したい
	// Redisのキャッシュキーを設定
	cacheKey := cache.GetIllustrationsListKey(int(req.Page))

	// Redisからキャッシュを取得
	var cachedResponse listIllustrationsResponse
	err := ctx.Server.RedisClient.Get(ctx.Context, cacheKey, &cachedResponse)
	if err != nil && !errors.Is(err, redis.Nil) {
		ctx.Server.Logger.Info("failed to redis err", zap.String("redis_key", cacheKey), zap.Error(err))
	}

	if err == nil {
		// キャッシュが存在する場合、それをレスポンスとして返す
		ctx.JSON(http.StatusOK, cachedResponse)
		return
	}

	arg := db.ListImageParams{
		Limit:  int32(ctx.Server.Config.ImageFetchLimit),
		Offset: int32(int(req.Page) * ctx.Server.Config.ImageFetchLimit),
	}
	images, err := ctx.Server.Store.ListImage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListImage : %w", err)))
		return
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := model.NewIllustration()
		il.Image = i

		illustrations = append(illustrations, il)
	}

	if len(illustrations) > 0 {
		response := listIllustrationsResponse{
			Illustrations: illustrations,
		}
		// Redisへのセットが失敗しても処理を続行
		err = ctx.Server.RedisClient.Set(ctx.Context, cacheKey, response, cache.CacheDurationDay)
		if err != nil {
			ctx.Server.Logger.Warn("failed redis data set", zap.String("redis_key", cacheKey), zap.Error(err))
		}
	}

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
	})
}

type getIllustrationsResponse struct {
	Illustration *model.Illustration `json:"illustration"`
}

// GetIllustration godoc
// @Summary Retrieve an illustration
// @Description Retrieves a single illustration by its ID
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID of the illustration to retrieve"
// @Success 200 {object} model/Illustration "The requested illustration"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Failed to parse 'id' number from path parameter"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No illustration found with the given ID"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to retrieve the illustration from the database"
// @Router /api/v1/illustrations/{id} [get]
func GetIllustration(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}

	// TODO: redis周りの処理は関数化したい
	// Redisのキャッシュキーを設定
	cacheKey := cache.GetIllustrationKey(id)

	// Redisからキャッシュを取得
	var cachedResponse getIllustrationsResponse
	err = ctx.Server.RedisClient.Get(ctx.Context, cacheKey, &cachedResponse)
	if err != nil && !errors.Is(err, redis.Nil) {
		// キャッシュの取得に失敗したが、デフォルトの動作としてDBからデータを取得する処理を続ける
		ctx.Server.Logger.Info("failed to redis err", zap.String("redis_key", cacheKey), zap.Error(err))
	}

	if err == nil {
		ctx.JSON(http.StatusOK, cachedResponse)
		return
	}

	image, err := ctx.Server.Store.GetImage(ctx.Context, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetImage: %w", err)))
			return
		}

		ctx.Server.Logger.Error("failed to GetImage", zap.Int("id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
		return
	}

	illustration := &model.Illustration{}
	illustration = service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, image)

	// レスポンスをキャッシュに保存
	response := getIllustrationsResponse{
		Illustration: illustration,
	}
	// Redisへのセットが失敗しても処理を続行
	err = ctx.Server.RedisClient.Set(ctx.Context, cacheKey, response, cache.CacheDurationWeek)
	if err != nil {
		ctx.Server.Logger.Warn("failed redis data set", zap.String("redis_key", cacheKey), zap.Error(err))
	}

	ctx.JSON(http.StatusOK, getIllustrationsResponse{
		Illustration: illustration,
	})
}

type searchIllustrationsRequest struct {
	Page  int    `form:"p"`
	Query string `form:"q"`
}

// TODO: imageだけでなく、カテゴリでも検索ができるようにする。
// また、検索結果をtrimし、被りがないようにする
// SearchIllustrations godoc
// @Summary Search illustrations
// @Description Searches for illustrations based on a query and page number.
// @Accept  json
// @Produce  json
// @Param   p     query   int    true  "Page number for pagination"
// @Param   q     query   string true  "Query string for searching illustrations by title or category"
// @Success 200   {array} model/Illustration "List of matched illustrations"
// @Failure 400   {object} request/JSONResponse{data=string} "Bad Request: The request is malformed or missing required fields."
// @Failure 500   {object} request/JSONResponse{data=string} "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// @Router /api/v1/illustrations/search [get]
func SearchIllustrations(ctx *app.AppContext) {
	var req searchIllustrationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	arg := db.SearchImagesParams{
		Limit:  int32(ctx.Server.Config.ImageFetchLimit),
		Offset: int32(req.Page * ctx.Server.Config.ImageFetchLimit),
		Query: sql.NullString{
			String: req.Query,
			Valid:  true,
		},
	}

	images, err := ctx.Server.Store.SearchImages(ctx, arg)
	if err != nil {
		ctx.Server.Logger.Error("failed to SearchImages", zap.String("query", req.Query), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to SearchImages : %w", err)))
		return
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := model.NewIllustration()
		il.Image = i

		illustrations = append(illustrations, il)
	}

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
	})
}

type listFetchRandomIllustrationsRequest struct {
	Limit       int64 `form:"limit"`
	ExclusionID int64 `form:"exclusion_id"`
}

// FetchRandomIllustrations godoc
// @Summary Fetch random illustrations
// @Description Retrieves a list of illustrations randomly selected from the database.
// @Accept  json
// @Produce  json
// @Param   limit   query   int  true  "Number of illustrations to retrieve"
// @Success 200 {array} models.Illustration "A list of illustrations"
// @Failure 400 {object} app.ErrorResponse "Bad Request: The request is malformed or missing required fields."
// @Failure 500 {object} app.ErrorResponse "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// @Router /api/v1/illustrations/random [get]
func FetchRandomIllustrations(ctx *app.AppContext) {
	var req listFetchRandomIllustrationsRequest
	if err := binder.BindQuery(ctx.Context, &req); err != nil {
		return
	}

	arg := db.FetchRandomImageParams{
		Limit: int32(req.Limit),
		ID:    req.ExclusionID,
	}
	images, err := ctx.Server.Store.FetchRandomImage(ctx, arg)
	if err != nil {
		ctx.Server.Logger.Error("failed to FetchRandomImage", zap.Int("exclusion_id", int(req.ExclusionID)), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to FetchRandomImage : %w", err)))
		return
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := model.NewIllustration()
		il.Image = i

		illustrations = append(illustrations, il)
	}

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
	})
}

type listIllustrationsByCharacterIDRequest struct {
	Page int64 `form:"p"`
}

// ListIllustrationsByCharacterID godoc
// @Summary List illustrations by character ID
// @Description Retrieves a paginated list of illustrations associated with a given character ID.
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID of the character"
// @Param   p    query  int  true  "Page number for pagination"
// @Success 200 {array} models.Illustration "A list of illustrations"
// @Failure 400 {object} app.ErrorResponse "Bad Request: The request is malformed or missing required fields."
// @Failure 404 {object} app.ErrorResponse "Not Found: No illustrations found for the given character ID."
// @Failure 500 {object} app.ErrorResponse "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// @Router /api/v1/illustrations/character/{id} [get]
func ListIllustrationsByCharacterID(ctx *app.AppContext) {
	charaID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}

	var req listIllustrationsByCharacterIDRequest
	if err := binder.BindQuery(ctx.Context, &req); err != nil {
		return
	}

	// TODO: redis周りの処理は関数化したい
	// Redisのキャッシュキーを設定
	cacheKey := cache.GetIllustrationsListByCharacterKey(charaID, int(req.Page))

	// Redisからキャッシュを取得
	var cachedResponse listIllustrationsResponse
	err = ctx.Server.RedisClient.Get(ctx.Context, cacheKey, &cachedResponse)
	if err != nil && !errors.Is(err, redis.Nil) {
		// キャッシュの取得に失敗したが、デフォルトの動作としてDBからデータを取得する処理を続ける
		ctx.Server.Logger.Info("failed to redis err", zap.String("redis_key", cacheKey), zap.Error(err))
	}

	if err == nil {
		// キャッシュが存在する場合、それをレスポンスとして返す
		ctx.JSON(http.StatusOK, cachedResponse)
		return
	}

	// charaIDを元にimage_characters_relationsを取得する
	arg := db.ListImageCharacterRelationsByCharacterIDWIthPaginationParams{
		Limit:       int32(ctx.Server.Config.ImageFetchLimit),
		Offset:      int32(int(req.Page) * ctx.Server.Config.ImageFetchLimit),
		CharacterID: int64(charaID),
	}
	icrs, err := ctx.Server.Store.ListImageCharacterRelationsByCharacterIDWIthPagination(ctx, arg)
	if err != nil {
		ctx.Server.Logger.Error("failed to ListImageCharacterRelationsByCharacterID", zap.Int("character_id", charaID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByCharacterID : %w", err)))
		return
	}

	// image_characters_relationsを元にimageを取得する
	images := []db.Image{}
	for _, icr := range icrs {
		image, err := ctx.Server.Store.GetImage(ctx, icr.ImageID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetImage: %w", err)))
				return
			}

			ctx.Server.Logger.Error("failed to GetImage", zap.Int64("image_id", icr.ImageID), zap.Int("character_id", charaID), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
			return
		}

		images = append(images, image)
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := model.NewIllustration()
		il.Image = i

		illustrations = append(illustrations, il)
	}

	if len(illustrations) > 0 {
		// レスポンスをキャッシュに保存
		response := listIllustrationsResponse{
			Illustrations: illustrations,
		}
		// Redisへのセットが失敗しても処理を続行
		err = ctx.Server.RedisClient.Set(ctx.Context, cacheKey, response, cache.CacheDurationDay)
		if err != nil {
			ctx.Server.Logger.Warn("failed redis data set", zap.String("redis_key", cacheKey), zap.Error(err))
		}
	}

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
	})
}

type listIllustrationsByChildCategoryIDRequest struct {
	Page int64 `form:"p"`
}

// ListIllustrationsByParentCategoryID godoc
// @Summary List illustrations by parent category ID
// @Description Retrieves a paginated list of illustrations associated with a given parent category ID.
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID of the parent category"
// @Param   p    query  int  true  "Page number for pagination"
// @Success 200 {array} models.Illustration "A list of illustrations"
// @Failure 400 {object} app.ErrorResponse "Bad Request: The request is malformed or missing required fields."
// @Failure 404 {object} app.ErrorResponse "Not Found: No illustrations found for the given parent category ID."
// @Failure 500 {object} app.ErrorResponse "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// @Router /api/v1/illustrations/category/child/{id} [get]
func ListIllustrationsByChildCategoryID(ctx *app.AppContext) {
	cCateID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}
	var req listIllustrationsByChildCategoryIDRequest
	if err := binder.BindQuery(ctx.Context, &req); err != nil {
		return
	}

	// TODO: redis周りの処理は関数化したい
	// Redisのキャッシュキーを設定
	cacheKey := cache.GetIllustrationsListByCategoryKey(cCateID, int(req.Page))

	// Redisからキャッシュを取得
	var cachedResponse listIllustrationsResponse
	err = ctx.Server.RedisClient.Get(ctx.Context, cacheKey, &cachedResponse)
	if err != nil && !errors.Is(err, redis.Nil) {
		// キャッシュの取得に失敗したが、デフォルトの動作としてDBからデータを取得する処理を続ける
		ctx.Server.Logger.Info("failed to redis err", zap.String("redis_key", cacheKey), zap.Error(err))
	}

	if err == nil {
		// キャッシュが存在する場合、それをレスポンスとして返す
		ctx.JSON(http.StatusOK, cachedResponse)
		return
	}

	// pCateIDを元にimage_parent_categories_relationsを取得する
	arg := db.ListImageChildCategoryRelationsByChildCategoryIDWithPaginationParams{
		Limit:           int32(ctx.Server.Config.ImageFetchLimit),
		Offset:          int32(int(req.Page) * ctx.Server.Config.ImageFetchLimit),
		ChildCategoryID: int64(cCateID),
	}
	icrs, err := ctx.Server.Store.ListImageChildCategoryRelationsByChildCategoryIDWithPagination(ctx, arg)
	if err != nil {
		ctx.Server.Logger.Error("failed to ListImageChildCategoryRelationsByChildCategoryIDWithPagination", zap.Int("child_category_id", cCateID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListImageChildCategoryRelationsByChildCategoryIDWithPagination : %w", err)))
		return
	}

	// image_parent_categories_relationsを元にimageを取得する
	images := []db.Image{}
	for _, icr := range icrs {
		image, err := ctx.Server.Store.GetImage(ctx, icr.ImageID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetImage: %w", err)))
				return
			}

			ctx.Server.Logger.Error("failed to GetImage", zap.Int64("image_id", icr.ImageID), zap.Int("child_category_id", cCateID), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
			return
		}

		images = append(images, image)
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := model.NewIllustration()
		il.Image = i

		illustrations = append(illustrations, il)
	}

	if len(illustrations) > 0 {
		// レスポンスをキャッシュに保存
		response := listIllustrationsResponse{
			Illustrations: illustrations,
		}
		// Redisへのセットが失敗しても処理を続行
		err = ctx.Server.RedisClient.Set(ctx.Context, cacheKey, response, cache.CacheDurationDay)
		if err != nil {
			ctx.Server.Logger.Warn("failed redis data set", zap.String("redis_key", cacheKey), zap.Error(err))
		}
	}

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
	})
}
