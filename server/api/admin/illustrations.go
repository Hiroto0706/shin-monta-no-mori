package admin

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"shin-monta-no-mori/internal/app"
	"shin-monta-no-mori/internal/cache"
	db "shin-monta-no-mori/internal/db/sqlc"
	model "shin-monta-no-mori/internal/domains/models"
	"shin-monta-no-mori/internal/domains/service"
	"shin-monta-no-mori/pkg/lib/binder"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	TotalPages    int64                 `json:"total_pages"`
	TotalCount    int64                 `json:"total_count"`
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
// @Router /api/v1/admin/illustrations/list [get]
func ListIllustrations(ctx *app.AppContext) {
	var req listIllustrationsRequest
	if err := binder.BindQuery(ctx.Context, &req); err != nil {
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
		ctx.Server.Logger.Error("failed to ListImage",
			zap.Int("offset", int(req.Page)),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListImage : %w", err)))
		return
	}

	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, i)

		illustrations = append(illustrations, il)
	}

	totalCount, err := ctx.Server.Store.CountImages(ctx)
	if err != nil {
		ctx.Server.Logger.Error("failed to CountImages",
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to CountImages : %w", err)))
		return
	}
	totalPages := (totalCount + int64(ctx.Server.Config.ImageFetchLimit-1)) / int64(ctx.Server.Config.ImageFetchLimit)

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
		TotalPages:    totalPages,
		TotalCount:    totalCount,
	})
}

type getIllustrationResponse struct {
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
// @Router /api/v1/admin/illustrations/{id} [get]
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

		ctx.Server.Logger.Error("failed to GetImage",
			zap.Int("illustration_id", id),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
		return
	}

	illustration := &model.Illustration{}
	illustration = service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, image)

	ctx.JSON(http.StatusOK, getIllustrationResponse{
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
// @Router /api/v1/admin/illustrations/search [get]
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
		ctx.Server.Logger.Error("failed to SearchImages",
			zap.String("query", req.Query),
			zap.Int("page", req.Page),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to SearchImages : %w", err)))
		return
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, i)

		illustrations = append(illustrations, il)
	}

	totalCount, err := ctx.Server.Store.CountSearchImages(ctx, sql.NullString{
		String: req.Query,
		Valid:  true,
	})
	if err != nil {
		ctx.Server.Logger.Error("failed to CountSearchImages",
			zap.String("query", req.Query),
			zap.Int("page", req.Page),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to CountSearchImages : %w", err)))
		return
	}
	totalPages := (totalCount + int64(ctx.Server.Config.ImageFetchLimit-1)) / int64(ctx.Server.Config.ImageFetchLimit)

	ctx.JSON(http.StatusOK, listIllustrationsResponse{
		Illustrations: illustrations,
		TotalPages:    totalPages,
		TotalCount:    totalCount,
	})
}

type createIllustrationRequest struct {
	Title             string               `form:"title" binding:"required"`
	Filename          string               `form:"filename" binding:"required"`
	Characters        []int64              `form:"characters[]"`
	ParentCategories  []int64              `form:"parent_categories[]"`
	ChildCategories   []int64              `form:"child_categories[]"`
	OriginalImageFile multipart.FileHeader `form:"original_image_file" binding:"required"`
	SimpleImageFile   multipart.FileHeader `form:"simple_image_file"`
}

// CreateIllustration godoc
// @Summary Create a new illustration
// @Description Creates a new illustration with title, filename, characters, categories, and image files.
// @Tags illustrations
// @Accept  multipart/form-data
// @Produce  json
// @Param   title              formData   string                 true  "Title of the illustration"
// @Param   filename           formData   string                 true  "Filename for the uploaded image"
// @Param   characters[]       formData   []int64                true  "List of character IDs associated with the illustration"
// @Param   parent_categories[] formData []int64                 true  "List of parent category IDs associated with the illustration"
// @Param   child_categories[] formData   []int64                true  "List of child category IDs associated with the illustration"
// @Param   original_image_file formData file                   true  "Original image file for the illustration"
// @Param   simple_image_file  formData  file                   false "Simple image file for the illustration (optional)"
// @Success 200 {object} gin/H "Returns the created illustration and a success message"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or validation"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to create the illustration due to a server error"
// @Router /api/v1/admin/illustrations/create [post]
func CreateIllustration(ctx *app.AppContext) {
	var req createIllustrationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to ShouldBind form data : %w", err)))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	var image db.Image
	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {

		var err error
		var originalSrc string
		if req.Filename != "" {
			originalSrc, err = service.UploadImageSrc(ctx.Context, &ctx.Server.Config, "original_image_file", req.Filename, IMAGE_TYPE_IMAGE, false)
			if err != nil {
				ctx.Server.Logger.Error("failed to UploadImage",
					zap.String("title", req.Title),
					zap.String("filename", req.Filename),
					zap.Error(err),
				)
				return fmt.Errorf("failed to UploadImage: %w", err)
			}
		}

		var simpleSrc string
		if req.Filename != "" && req.SimpleImageFile.Size != 0 {
			simpleSrc, err = service.UploadImageSrc(ctx.Context, &ctx.Server.Config, "simple_image_file", req.Filename, IMAGE_TYPE_IMAGE, true)
			if err != nil {
				ctx.Server.Logger.Error("failed to UploadImage for simpl image",
					zap.String("title", req.Title),
					zap.String("filename", req.Filename),
					zap.Error(err),
				)
				return fmt.Errorf("failed to UploadImage for simpl image : %w", err)
			}
		}

		arg := db.CreateImageParams{
			Title:            req.Title,
			OriginalSrc:      originalSrc,
			OriginalFilename: req.Filename,
			SimpleSrc:        sql.NullString{String: "", Valid: false},
			SimpleFilename:   sql.NullString{String: "", Valid: false},
		}
		if simpleSrc != "" {
			arg.SimpleSrc = sql.NullString{String: simpleSrc, Valid: true}
			arg.SimpleFilename = sql.NullString{String: req.Filename + "_s", Valid: true}
		}

		image, err = ctx.Server.Store.CreateImage(ctx, arg)
		if err != nil {
			ctx.Server.Logger.Error("failed to CreateImage",
				zap.String("title", req.Title),
				zap.String("filename", req.Filename),
				zap.String("original_src", originalSrc),
				zap.String("simple_src", simpleSrc),
				zap.Error(err),
			)
			return fmt.Errorf("failed to CreateImage: %w", err)
		}

		// ImageCharacterRelationsの保存
		for _, c_id := range req.Characters {
			arg := db.CreateImageCharacterRelationsParams{
				ImageID:     image.ID,
				CharacterID: c_id,
			}
			_, err := ctx.Server.Store.CreateImageCharacterRelations(ctx, arg)
			if err != nil {
				ctx.Server.Logger.Error("failed to CreateImageCharacterRelations",
					zap.String("title", req.Title),
					zap.String("filename", req.Filename),
					zap.String("original_src", originalSrc),
					zap.String("simple_src", simpleSrc),
					zap.Int("character_id", int(c_id)),
					zap.Error(err),
				)
				return fmt.Errorf("failed to CreateImageCharacterRelations: %w", err)
			}
		}

		// ImageParentCategoryRelationsの保存
		// mapを使うことで、重複する値を取り除く
		parentCategorySet := make(map[int64]struct{})

		for _, pc_id := range req.ParentCategories {
			parentCategorySet[pc_id] = struct{}{}
		}

		for pc_id := range parentCategorySet {
			arg := db.CreateImageParentCategoryRelationsParams{
				ImageID:          image.ID,
				ParentCategoryID: pc_id,
			}

			_, err := ctx.Server.Store.CreateImageParentCategoryRelations(ctx, arg)
			if err != nil {
				ctx.Server.Logger.Error("failed to CreateImageParentCategoryRelations",
					zap.String("title", req.Title),
					zap.String("filename", req.Filename),
					zap.String("original_src", originalSrc),
					zap.String("simple_src", simpleSrc),
					zap.Int("parent_category_id", int(pc_id)),
					zap.Error(err),
				)
				return fmt.Errorf("failed to CreateImageParentCategoryRelations: %w", err)
			}
		}

		// ImageChildCategoryRelationsの保存
		for _, cc_id := range req.ChildCategories {
			arg := db.CreateImageChildCategoryRelationsParams{
				ImageID:         image.ID,
				ChildCategoryID: cc_id,
			}
			_, err := ctx.Server.Store.CreateImageChildCategoryRelations(ctx, arg)
			if err != nil {
				ctx.Server.Logger.Error("failed to CreateImageChildCategoryRelations",
					zap.String("title", req.Title),
					zap.String("filename", req.Filename),
					zap.String("original_src", originalSrc),
					zap.String("simple_src", simpleSrc),
					zap.Int("child_category_id", int(cc_id)),
					zap.Error(err),
				)
				return fmt.Errorf("failed to CreateImageChildCategoryRelations: %w", err)
			}
		}

		return nil
	})

	if txErr != nil {
		ctx.Server.Logger.Error("CreateImage transaction was failed",
			zap.String("title", req.Title),
			zap.String("filename", req.Filename),
			zap.Error(txErr),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("CreateImage transaction was failed : %w", txErr)))
		return
	}

	// イラスト一覧を最新にするためにillustrationsを取得
	image, err := ctx.Server.Store.GetImage(ctx, int64(image.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
			return
		}

		ctx.Server.Logger.Error("failed to GetImage",
			zap.Int("illustration_id", int(image.ID)),
			zap.String("title", req.Title),
			zap.String("filename", req.Filename),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
		return
	}

	illustration := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, image)

	// redisキャッシュの削除
	keyPattern := []string{cache.IllustrationsPrefix + "*"}
	err = ctx.Server.RedisClient.Del(ctx, keyPattern)
	if err != nil {
		ctx.Server.Logger.Warn("failed redis data delete", zap.Error(err))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"illustration": illustration,
		"message":      "illustrationの作成に成功しました",
	})
}

type editIllustrationRequest struct {
	Title               string               `form:"title"`
	Filename            string               `form:"filename"`
	Characters          []int64              `form:"characters[]"`
	ParentCategories    []int64              `form:"parent_categories[]"`
	ChildCategories     []int64              `form:"child_categories[]"`
	OriginalImageFile   multipart.FileHeader `form:"original_image_file"`
	SimpleImageFile     multipart.FileHeader `form:"simple_image_file"`
	IsDeleteSimpleImage bool                 `form:"is_delete_simple_image"`
}

// EditIllustration godoc
// @Summary Edit an illustration
// @Description Updates an illustration by its ID with new title, filename, and optionally updates the image.
// @Accept  multipart/form-data
// @Produce  json
// @Param   id          path     int    true  "ID of the illustration to update"
// @Param   title       formData string true  "New title of the illustration"
// @Param   filename    formData string true  "New filename for the illustration; used in image re-upload"
// @Param   image_file  formData file   false "New image file for the illustration"
// @Param   characters  formData []int  false "List of character IDs associated with the illustration"
// @Param   parentCategories formData []int false "List of parent category IDs associated with the illustration"
// @Param   childCategories  formData []int false "List of child category IDs associated with the illustration"
// @Success 200 {object} gin/H "Returns the updated illustration and a success message"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or validation"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No illustration found with the given ID"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to update the illustration due to a server error"
// @Router /api/v1/admin/illustrations/{id} [put]
func EditIllustration(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	var req editIllustrationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	image, err := ctx.Server.Store.GetImage(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
			return
		}
		ctx.Server.Logger.Error("failed to GetImage",
			zap.Int("illustration_id", id),
			zap.String("title", req.Title),
			zap.String("filename", req.Filename),
			zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
		return
	}

	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {
		// Conditions for updating simpleSrc:
		// 1. ファイル名のみ変更
		// 2. イメージのみ変更
		// 3. ファイル名＆イメージが変更
		originalSrc := image.OriginalSrc
		if image.OriginalFilename != req.Filename || req.OriginalImageFile.Filename != "" {
			err := service.DeleteImageSrc(ctx.Context, &ctx.Server.Config, image.OriginalSrc)
			if err != nil {
				ctx.Server.Logger.Error("failed to DeleteImageSrc",
					zap.Int("illustration_id", id),
					zap.String("title", req.Title),
					zap.String("filename", req.Filename),
					zap.String("original_src", originalSrc),
					zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
					zap.Error(err),
				)
				return fmt.Errorf("failed to DeleteImageSrc : %w", err)
			}

			originalSrc, err = service.UploadImageSrc(ctx.Context, &ctx.Server.Config, "original_image_file", req.Filename, IMAGE_TYPE_IMAGE, false)
			if err != nil {
				ctx.Server.Logger.Error("failed to UploadImage",
					zap.Int("illustration_id", id),
					zap.String("title", req.Title),
					zap.String("filename", req.Filename),
					zap.String("original_src", originalSrc),
					zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
					zap.Error(err),
				)
				return fmt.Errorf("failed to UploadImage: %w", err)
			}
		}

		// Conditions for updating simpleSrc:
		// 1. ファイル名のみ変更
		// 2. イメージのみ変更
		// 3. ファイル名＆イメージが変更
		// 4. イメージの削除
		shouldUpdateSimpleSrc := image.OriginalFilename != req.Filename || req.SimpleImageFile.Filename != ""
		simpleSrc := image.SimpleSrc.String
		if shouldUpdateSimpleSrc {
			if simpleSrc != "" {
				err := service.DeleteImageSrc(ctx.Context, &ctx.Server.Config, simpleSrc)
				if err != nil {
					ctx.Server.Logger.Error("failed to DeleteImageSrc fro simple image",
						zap.Int("illustration_id", id),
						zap.String("title", req.Title),
						zap.String("filename", req.Filename),
						zap.String("simple_src", simpleSrc),
						zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
						zap.Error(err),
					)
					return fmt.Errorf("failed to DeleteImageSrc fro simple : %w", err)
				}
			}

			if req.SimpleImageFile.Filename != "" {
				simpleSrc, err = service.UploadImageSrc(ctx.Context, &ctx.Server.Config, "simple_image_file", req.Filename, IMAGE_TYPE_IMAGE, true)
				if err != nil {
					ctx.Server.Logger.Error("failed to UploadImage fro simple image",
						zap.Int("illustration_id", id),
						zap.String("title", req.Title),
						zap.String("filename", req.Filename),
						zap.String("simple_src", simpleSrc),
						zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
						zap.Error(err),
					)
					return fmt.Errorf("failed to UploadImage for simple image : %w", err)
				}
			}
		}

		if req.IsDeleteSimpleImage {
			err := service.DeleteImageSrc(ctx.Context, &ctx.Server.Config, simpleSrc)
			if err != nil {
				ctx.Server.Logger.Error("failed to DeleteImageSrc fro simple image",
					zap.Int("illustration_id", id),
					zap.String("title", req.Title),
					zap.String("filename", req.Filename),
					zap.String("simple_src", simpleSrc),
					zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
					zap.Error(err),
				)
				return fmt.Errorf("failed to DeleteImageSrc for simple image : %w", err)
			}
			simpleSrc = ""
		}

		// imageのUpdate処理
		arg := db.UpdateImageParams{
			ID:               image.ID,
			Title:            req.Title,
			OriginalSrc:      originalSrc,
			SimpleSrc:        sql.NullString{String: "", Valid: false},
			OriginalFilename: req.Filename,
			SimpleFilename:   sql.NullString{String: "", Valid: false},
			// TODO: timezoneがUTCになっている。厳密な時系列を扱う必要がある課題が出た時に修正する必要あり。
			UpdatedAt: time.Now(),
		}
		if simpleSrc != "" {
			arg.SimpleSrc = sql.NullString{String: simpleSrc, Valid: true}
			arg.SimpleFilename = sql.NullString{String: req.Filename + "_s", Valid: true}
		}
		image, err = ctx.Server.Store.UpdateImage(ctx, arg)
		if err != nil {
			ctx.Server.Logger.Error("failed to UpdateImage",
				zap.Int("illustration_id", id),
				zap.String("title", req.Title),
				zap.String("filename", req.Filename),
				zap.String("original_src", originalSrc),
				zap.String("simple_src", simpleSrc),
				zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
				zap.Error(err),
			)
			return err
		}

		// TODO: relation周りのUpdate処理は共通化できそう
		// image_character_relationsのUpdate処理
		err = service.UpdateImageCharacterRelationsIDs(ctx.Context, ctx.Server.Store, image.ID, req.Characters)
		if err != nil {
			ctx.Server.Logger.Error("failed to UpdateImageCharacterRelationsIDs",
				zap.Int("illustration_id", id),
				zap.String("title", req.Title),
				zap.String("filename", req.Filename),
				zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
				zap.Error(err),
			)
			return fmt.Errorf("failed to server.UpdateImageCharacterRelationsIDs: %w", err)
		}

		// image_parent_category_relationsのUpdate処理
		err = service.UpdateImageParentCategoryRelationsIDs(ctx.Context, ctx.Server.Store, image.ID, req.ParentCategories)
		if err != nil {
			ctx.Server.Logger.Error("failed to UpdateImageParentCategoryRelationsIDs",
				zap.Int("illustration_id", id),
				zap.String("title", req.Title),
				zap.String("filename", req.Filename),
				zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
				zap.Error(err),
			)
			return fmt.Errorf("failed to server.UpdateImageParentCategoryRelationsIDs: %w", err)
		}

		// image_child_category_relationsのUpdate処理
		err = service.UpdateImageChildCategoryRelationsIDs(ctx.Context, ctx.Server.Store, image.ID, req.ChildCategories)
		if err != nil {
			ctx.Server.Logger.Error("failed to UpdateImageChildCategoryRelationsIDs",
				zap.Int("illustration_id", id),
				zap.String("title", req.Title),
				zap.String("filename", req.Filename),
				zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
				zap.Error(err),
			)
			return fmt.Errorf("failed to server.UpdateImageChildCategoryRelationsIDs: %w", err)
		}

		return nil
	})

	if txErr != nil {
		ctx.Server.Logger.Error("EditImage transaction was failed",
			zap.Int("illustration_id", id),
			zap.String("title", req.Title),
			zap.String("filename", req.Filename),
			zap.Bool("is_delete_simple_image", req.IsDeleteSimpleImage),
			zap.Error(txErr),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(txErr))
		return
	}

	// redisキャッシュの削除
	keyPattern := []string{
		cache.IllustrationsPrefix + "*",
		cache.GetIllustrationKey(id),
	}
	err = ctx.Server.RedisClient.Del(ctx, keyPattern)
	if err != nil {
		ctx.Server.Logger.Warn("failed redis data delete", zap.Error(err))
	}

	illustration := service.FetchRelationInfoForIllustrations(ctx.Context, ctx.Server.Store, image)

	ctx.JSON(http.StatusOK, gin.H{
		"illustration": illustration,
		"message":      "illustrationの編集に成功しました",
	})
}

// DeleteIllustration godoc
// @Summary Delete an illustration
// @Description Deletes a specific illustration by its ID.
// @Tags illustrations
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID of the illustration to delete"
// @Success 200 {object} gin/H "Returns a success message indicating the illustration has been deleted"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error parsing the 'id' from path parameters"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No illustration found with the given ID"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to delete the illustration due to a server error"
// @Router /api/v1/admin/illustrations/{id} [delete]
func DeleteIllustration(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	// TODO: illustrationとして取得すれば、冗長な関数を削除できそう
	image, err := ctx.Server.Store.GetImage(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
			return
		}
		ctx.Server.Logger.Error("failed to GetImage",
			zap.Int("illustration_id", id),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
		return
	}

	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {
		err = service.DeleteImageSrc(ctx.Context, &ctx.Server.Config, image.OriginalSrc)
		if err != nil {
			ctx.Server.Logger.Error("failed to DeleteImageSrc",
				zap.Int("illustration_id", id),
				zap.String("original_src", image.OriginalSrc),
				zap.Error(err),
			)
			return fmt.Errorf("failed to DeleteImageSrc : %w", err)
		}

		if image.SimpleSrc.String != "" {
			err = service.DeleteImageSrc(ctx.Context, &ctx.Server.Config, image.SimpleSrc.String)
			if err != nil {
				ctx.Server.Logger.Error("failed to DeleteImageSrc for simple image",
					zap.Int("illustration_id", id),
					zap.String("simple_src", image.SimpleFilename.String),
					zap.Error(err),
				)
				return fmt.Errorf("failed to DeleteImageSrc for simple image : %w", err)
			}
		}

		// TODO: illustrationとして取得できれば、このrelation取得の処理削除できる
		// image_child_category_relationsを削除
		err = ctx.Server.Store.DeleteAllImageChildCategoryRelationsByImageID(ctx, image.ID)
		if err != nil {
			ctx.Server.Logger.Error("failed to DeleteAllImageChildCategoryRelationsByImageID",
				zap.Int("illustration_id", id),
				zap.Error(err),
			)
			return fmt.Errorf("failed to DeleteAllImageChildCategoryRelationsByImageID: %w", err)
		}

		// image_parent_category_relationsを削除
		err = ctx.Server.Store.DeleteAllImageParentCategoryRelationsByImageID(ctx, image.ID)
		if err != nil {
			ctx.Server.Logger.Error("failed to DeleteAllImageParentCategoryRelationsByImageID",
				zap.Int("illustration_id", id),
				zap.Error(err),
			)
			return fmt.Errorf("failed to DeleteAllImageParentCategoryRelationsByImageID: %w", err)
		}

		// image_character_relationsを削除
		err = ctx.Server.Store.DeleteAllImageCharacterRelationsByImageID(ctx, image.ID)
		if err != nil {
			ctx.Server.Logger.Error("failed to DeleteAllImageCharacterRelationsByImageID",
				zap.Int("illustration_id", id),
				zap.Error(err),
			)
			return fmt.Errorf("failed to DeleteAllImageCharacterRelationsByImageID: %w", err)
		}

		// Imageを削除
		err = ctx.Server.Store.DeleteImage(ctx, image.ID)
		if err != nil {
			ctx.Server.Logger.Error("failed to DeleteImage",
				zap.Int("illustration_id", id),
				zap.Error(err),
			)
			return fmt.Errorf("failed to DeleteImage: %w", err)
		}

		return nil
	})

	if txErr != nil {
		ctx.Server.Logger.Error("DeleteImage transaction was failed",
			zap.Int("illustration_id", id),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(txErr))
		return
	}

	// redisキャッシュの削除
	keyPattern := []string{
		cache.IllustrationsPrefix + "*",
		cache.GetIllustrationKey(id),
	}
	err = ctx.Server.RedisClient.Del(ctx, keyPattern)
	if err != nil {
		ctx.Server.Logger.Warn("failed redis data delete", zap.Error(err))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "illustrationの削除に成功しました",
	})
}
