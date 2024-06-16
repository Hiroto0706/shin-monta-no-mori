package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/internal/domains/service"
	"strconv"
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
	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listIllustrationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	illustrations := []*model.Illustration{}

	arg := db.ListImageParams{
		Limit:  int32(ctx.Server.Config.ImageFetchLimit),
		Offset: int32(int(req.Page) * ctx.Server.Config.ImageFetchLimit),
	}
	images, err := ctx.Server.Store.ListImage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListImage : %w", err)))
		return
	}

	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, i)

		illustrations = append(illustrations, il)
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

	image, err := ctx.Server.Store.GetImage(ctx.Context, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetImage: %w", err)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
		return
	}

	illustration := &model.Illustration{}
	illustration = service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, image)

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
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to SearchImages : %w", err)))
		return
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, i)

		illustrations = append(illustrations, il)
	}

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
	})
}

type listFetchRandomIllustrationsRequest struct {
	Limit int64 `form:"limit"`
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
	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listFetchRandomIllustrationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	illustrations := []*model.Illustration{}

	images, err := ctx.Server.Store.FetchRandomImage(ctx, int32(req.Limit))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to FetchRandomImage : %w", err)))
		return
	}

	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, i)

		illustrations = append(illustrations, il)
	}

	ctx.JSON(http.StatusOK, illustrations)
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
	log.Println("ここきてる？")
	charaID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}

	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listIllustrationsByCharacterIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
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

			ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
			return
		}

		images = append(images, image)
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, i)

		illustrations = append(illustrations, il)
	}

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
	})
}

type listIllustrationsByParentCategoryIDRequest struct {
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
// @Router /api/v1/illustrations/category/parent/{id} [get]
func ListIllustrationsByParentCategoryID(ctx *app.AppContext) {
	pCateID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}

	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listIllustrationsByParentCategoryIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	// pCateIDを元にimage_parent_categories_relationsを取得する
	arg := db.ListImageParentCategoryRelationsByParentCategoryIDWithPaginationParams{
		Limit:            int32(ctx.Server.Config.ImageFetchLimit),
		Offset:           int32(int(req.Page) * ctx.Server.Config.ImageFetchLimit),
		ParentCategoryID: int64(pCateID),
	}
	icrs, err := ctx.Server.Store.ListImageParentCategoryRelationsByParentCategoryIDWithPagination(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListImageParentCategoryRelationsByParentCategoryIDWithPagination : %w", err)))
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

			ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
			return
		}

		images = append(images, image)
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, i)

		illustrations = append(illustrations, il)
	}

	ctx.JSON(http.StatusOK, illustrations)
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

	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listIllustrationsByChildCategoryIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
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

			ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
			return
		}

		images = append(images, image)
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, i)

		illustrations = append(illustrations, il)
	}

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
	})
}
