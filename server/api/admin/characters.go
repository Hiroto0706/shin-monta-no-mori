package admin

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"shin-monta-no-mori/internal/app"
	"shin-monta-no-mori/internal/cache"
	db "shin-monta-no-mori/internal/db/sqlc"
	"shin-monta-no-mori/internal/domains/service"
	"shin-monta-no-mori/pkg/lib/binder"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	IMAGE_TYPE_CHARACTER = "character"
)

type listCharactersRequest struct {
	Page int64 `form:"p"`
}

type listCharactersResponse struct {
	Characters []db.Character `json:"characters"`
	TotalPages int64          `json:"total_pages"`
	TotalCount int64          `json:"total_count"`
}

// ListCharacters godoc
// @Summary List characters
// @Description Retrieves a paginated list of characters based on the provided page number.
// @Accept  json
// @Produce  json
// @Param   p     query   int64  true  "Page number for pagination"
// @Success 200   {object} gin/H  "Returns a list of characters"
// @Failure 400   {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or validation"
// @Failure 500   {object} request/JSONResponse{data=string} "Internal Server Error: Failed to list the characters"
// @Router /api/v1/admin/characters/list [get]
func ListCharacters(ctx *app.AppContext) {
	var req listCharactersRequest
	if err := binder.BindQuery(ctx.Context, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	arg := db.ListCharactersParams{
		Limit:  int32(ctx.Server.Config.CharacterFetchLimit),
		Offset: int32(int(req.Page) * ctx.Server.Config.CharacterFetchLimit),
	}
	characters, err := ctx.Server.Store.ListCharacters(ctx, arg)
	if err != nil {
		ctx.Server.Logger.Error("failed to ListAllCharacters", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListAllCharacters : %w", err)))
		return
	}

	totalCount, err := ctx.Server.Store.CountCharacters(ctx)
	if err != nil {
		ctx.Server.Logger.Error("failed to CountCharacters", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to CountCharacters : %w", err)))
		return
	}
	totalPages := (totalCount + int64(ctx.Server.Config.CharacterFetchLimit-1)) / int64(ctx.Server.Config.CharacterFetchLimit)

	ctx.JSON(http.StatusOK, listCharactersResponse{
		Characters: characters,
		TotalPages: totalPages,
		TotalCount: totalCount,
	})
}

type listAllCharactersResponse struct {
	Characters []db.Character `json:"characters"`
}

// ListAllCharacters godoc
// @Summary List characters
// @Description Retrieves a paginated list of characters based on the provided page number.
// @Accept  json
// @Produce  json
// @Param   p     query   int64  true  "Page number for pagination"
// @Success 200   {object} gin/H  "Returns a list of characters"
// @Failure 400   {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or validation"
// @Failure 500   {object} request/JSONResponse{data=string} "Internal Server Error: Failed to list the characters"
// @Router /api/v1/admin/characters/list/all [get]
func ListAllCharacters(ctx *app.AppContext) {
	characters, err := ctx.Server.Store.ListAllCharacters(ctx)
	if err != nil {
		ctx.Server.Logger.Error("failed to ListAllCharacters", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListAllCharacters : %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, listAllCharactersResponse{
		Characters: characters,
	})
}

type searchCharactersRequest struct {
	Page  int    `form:"p"`
	Query string `form:"q"`
}

// SearchCharacters godoc
// @Summary Search characters
// @Description Searches for characters based on a query string and page number.
// @Accept  json
// @Produce  json
// @Param   p     query   int    true  "Page number for pagination"
// @Param   q     query   string true  "Query string for searching characters by name or other attributes"
// @Success 200   {object} gin/H  "Returns a list of characters that match the query"
// @Failure 400   {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or validation"
// @Failure 500   {object} request/JSONResponse{data=string} "Internal Server Error: Failed to search the characters"
// @Router /api/v1/admin/characters/search [get]
func SearchCharacters(ctx *app.AppContext) {
	var req searchCharactersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	arg := db.SearchCharactersParams{
		Limit:  int32(ctx.Server.Config.CharacterFetchLimit),
		Offset: int32(req.Page * ctx.Server.Config.CharacterFetchLimit),
		Query: sql.NullString{
			String: req.Query,
			Valid:  true,
		},
	}
	characters, err := ctx.Server.Store.SearchCharacters(ctx, arg)
	if err != nil {
		ctx.Server.Logger.Error("failed to SearchCharacters", zap.Int("page", req.Page), zap.String("query", req.Query), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to SearchCharacters : %w", err)))
		return
	}

	totalCount, err := ctx.Server.Store.CountSearchCharacters(ctx, sql.NullString{
		String: req.Query,
		Valid:  true,
	})
	if err != nil {
		ctx.Server.Logger.Error("failed to CountSearchCharacters", zap.Int("page", req.Page), zap.String("query", req.Query), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to CountSearchCharacters : %w", err)))
		return
	}
	totalPages := (totalCount + int64(ctx.Server.Config.CharacterFetchLimit-1)) / int64(ctx.Server.Config.CharacterFetchLimit)

	ctx.JSON(http.StatusOK, listCharactersResponse{
		Characters: characters,
		TotalPages: totalPages,
		TotalCount: totalCount,
	})
}

type getCharacterResponse struct {
	Character db.Character `json:"character"`
}

// GetCharacter godoc
// @Summary Retrieve a character
// @Description Retrieves a single character by its ID.
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID of the character to retrieve"
// @Success 200 {object} gin/H "The requested character"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Failed to parse 'id' number from path parameter"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No character found with the given ID"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to retrieve the character from the database"
// @Router /api/v1/admin/characters/{id} [get]
func GetCharacter(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}

	character, err := ctx.Server.Store.GetCharacter(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetCharacter: %w", err)))
			return
		}

		ctx.Server.Logger.Error("failed to GetCharacter", zap.Int("character_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, getCharacterResponse{
		Character: character,
	})
}

type createCharacterRequest struct {
	Name          string               `form:"name" binding:"required"`
	Filename      string               `form:"filename" binding:"required"`
	ImageFile     multipart.FileHeader `form:"image_file" binding:"required"`
	PriorityLevel int16                `form:"priority_level" binding:"required"`
}

// CreateCharacter godoc
// @Summary Create a new character
// @Description Creates a new character with a name, filename, and image file.
// @Accept  multipart/form-data
// @Produce  json
// @Param   name       formData   string  true  "Name of the character"
// @Param   filename   formData   string  true  "Filename for the uploaded image"
// @Param   image_file formData   file    true  "Image file for the character"
// @Success 200 {object} gin/H "Returns the created character and a success message"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or validation"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to create the character due to a transaction error"
// @Router /api/v1/admin/characters/create [post]
func CreateCharacter(ctx *app.AppContext) {
	var req createCharacterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	var character db.Character
	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {
		var src string
		src, err := service.UploadImageSrc(ctx.Context, &ctx.Server.Config, "image_file", req.Filename, IMAGE_TYPE_CHARACTER, false)
		if err != nil {
			ctx.Server.Logger.Error("failed to UploadImage", zap.String("name", req.Name), zap.Int16("name", req.PriorityLevel), zap.Error(err))
			return fmt.Errorf("failed to UploadImage: %w", err)
		}

		arg := db.CreateCharacterParams{
			Name:          req.Name,
			Src:           src,
			Filename:      sql.NullString{String: req.Filename, Valid: true},
			PriorityLevel: req.PriorityLevel,
		}
		character, err = ctx.Server.Store.CreateCharacter(ctx, arg)
		if err != nil {
			ctx.Server.Logger.Error("failed to CreateCharacter", zap.String("name", req.Name), zap.Int16("name", req.PriorityLevel), zap.Error(err))
			return fmt.Errorf("failed to CreateCharacter: %w", err)
		}

		return nil
	})

	if txErr != nil {
		ctx.Server.Logger.Error("CreateCharacter transaction was failed", zap.Error(txErr))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("CreateCharacter transaction was failed : %w", txErr)))
		return
	}

	// redisキャッシュの削除
	keyPattern := []string{cache.CharactersPrefix + "*"}
	err := ctx.Server.RedisClient.Del(ctx, keyPattern)
	if err != nil {
		ctx.Server.Logger.Warn("failed redis data delete", zap.Error(err))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"character": character,
		"message":   "キャラクターの作成に成功しました",
	})
}

type editCharacterRequest struct {
	Name          string               `form:"name"`
	Filename      string               `form:"filename"`
	ImageFile     multipart.FileHeader `form:"image_file"`
	PriorityLevel int16                `form:"priority_level"`
}

// EditCharacter godoc
// @Summary Edit an existing character
// @Description Edits an existing character by its ID, updating its name, filename, and image file.
// @Accept  multipart/form-data
// @Produce  json
// @Param   id         path     int                   true  "ID of the character to edit"
// @Param   name       formData string                true  "New name for the character"
// @Param   filename   formData string                true  "New filename for the uploaded image"
// @Param   image_file formData file                  true  "New image file for the character"
// @Success 200        {object} gin/H                 "Returns the updated character and a success message"
// @Failure 400        {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or path parameter parsing"
// @Failure 404        {object} request/JSONResponse{data=string} "Not Found: No character found with the given ID"
// @Failure 500        {object} request/JSONResponse{data=string} "Internal Server Error: Failed to edit the character due to a transaction error"
// @Router /api/v1/admin/characters/{id} [put]
func EditCharacter(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}
	var req editCharacterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	character, err := ctx.Server.Store.GetCharacter(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
			return
		}
		ctx.Server.Logger.Error("failed to GetCharacter", zap.Int("character_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
		return
	}

	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {
		src := character.Src
		if character.Filename.String != req.Filename || req.ImageFile.Filename != "" {
			err := service.DeleteImageSrc(ctx.Context, &ctx.Server.Config, character.Src)
			if err != nil {
				ctx.Server.Logger.Error("failed to DeleteImageSrc", zap.Int("character_id", id), zap.Error(err))
				return fmt.Errorf("failed to DeleteImageSrc : %w", err)
			}

			src, err = service.UploadImageSrc(ctx.Context, &ctx.Server.Config, "image_file", req.Filename, IMAGE_TYPE_CHARACTER, false)
			if err != nil {
				ctx.Server.Logger.Error("failed to UploadImage", zap.Int("character_id", id), zap.Error(err))
				return fmt.Errorf("failed to UploadImage : %w", err)
			}
		}

		arg := db.UpdateCharacterParams{
			ID:            character.ID,
			Name:          req.Name,
			Src:           src,
			Filename:      sql.NullString{String: character.Filename.String, Valid: true},
			PriorityLevel: req.PriorityLevel,
			UpdatedAt:     time.Now(),
		}
		if character.Filename.String != req.Filename {
			arg.Filename = sql.NullString{String: req.Filename, Valid: true}
		}

		character, err = ctx.Server.Store.UpdateCharacter(ctx, arg)
		if err != nil {
			ctx.Server.Logger.Error("failed to UpdateCharacter", zap.Int("character_id", id), zap.Error(err))
			return err
		}

		return nil
	})

	if txErr != nil {
		ctx.Server.Logger.Error("EditCharacter transaction was failed", zap.Int("character_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("EditCharacter transaction was failed : %w", txErr)))
		return
	}

	// redisキャッシュの削除
	keyPattern := []string{cache.CharactersPrefix + "*"}
	err = ctx.Server.RedisClient.Del(ctx, keyPattern)
	if err != nil {
		ctx.Server.Logger.Warn("failed redis data delete", zap.Error(err))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"character": character,
		"message":   "characterの編集に成功しました",
	})
}

// DeleteCharacter godoc
// @Summary Delete a character
// @Description Deletes a character by its ID along with associated resources.
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID of the character to delete"
// @Success 200   {object} gin/H "Returns a success message upon successful deletion"
// @Failure 400   {object} request/JSONResponse{data=string} "Bad Request: Error parsing character ID from path parameter"
// @Failure 404   {object} request/JSONResponse{data=string} "Not Found: No character found with the given ID"
// @Failure 500   {object} request/JSONResponse{data=string} "Internal Server Error: Failed to delete the character due to a transaction error"
// @Router /api/v1/admin/characters/{id} [delete]
func DeleteCharacter(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	character, err := ctx.Server.Store.GetCharacter(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
			return
		}
		ctx.Server.Logger.Error("failed to GetCharacter", zap.Int("character_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
		return
	}

	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {
		err = service.DeleteImageSrc(ctx.Context, &ctx.Server.Config, character.Src)
		if err != nil {
			ctx.Server.Logger.Error("failed to DeleteImageSrc", zap.Int("character_id", id), zap.Error(err))
			return fmt.Errorf("failed to DeleteImageSrc: %w", err)
		}

		// images_character_relationsの削除
		err = ctx.Server.Store.DeleteAllImageCharacterRelationsByCharacterID(ctx, character.ID)
		if err != nil {
			ctx.Server.Logger.Error("failed to DeleteAllImageCharacterRelationsByCharacterID", zap.Int("character_id", id), zap.Error(err))
			return fmt.Errorf("failed to DeleteAllImageCharacterRelationsByCharacterID : %w", err)
		}

		err = ctx.Server.Store.DeleteCharacter(ctx, int64(id))
		if err != nil {
			ctx.Server.Logger.Error("failed to DeleteCharacter", zap.Int("character_id", id), zap.Error(err))
			return fmt.Errorf("failed to DeleteCharacter: %w", err)
		}
		return nil
	})

	if txErr != nil {
		ctx.Server.Logger.Error("DeleteCharacter transaction was failed", zap.Int("character_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("DeleteCharacter transaction was failed : %w", txErr)))
		return
	}

	// redisキャッシュの削除
	keyPattern := []string{cache.CharactersPrefix + "*"}
	err = ctx.Server.RedisClient.Del(ctx, keyPattern)
	if err != nil {
		ctx.Server.Logger.Warn("failed redis data delete", zap.Error(err))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "characterの削除に成功しました",
	})
}
