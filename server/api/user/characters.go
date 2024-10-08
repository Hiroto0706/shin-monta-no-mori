package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"shin-monta-no-mori/internal/app"
	"shin-monta-no-mori/internal/cache"
	db "shin-monta-no-mori/internal/db/sqlc"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	IMAGE_TYPE_CHARACTER = "character"
)

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
// @Router /api/v1/characters/list/all [get]
func ListAllCharacters(ctx *app.AppContext) {
	// TODO: redis周りの処理は関数化したい
	// Redisのキャッシュキーを設定
	cacheKey := cache.GetCharactersAllKey()

	// Redisからキャッシュを取得
	var cachedResponse listAllCharactersResponse
	err := ctx.Server.RedisClient.Get(ctx.Context, cacheKey, &cachedResponse)
	if err != nil && !errors.Is(err, redis.Nil) {
		// キャッシュの取得に失敗したが、デフォルトの動作としてDBからデータを取得する処理を続ける
		ctx.Server.Logger.Info("failed to redis err", zap.String("redis_key", cacheKey), zap.Error(err))
	}

	if err == nil {
		// キャッシュが存在する場合、それをレスポンスとして返す
		ctx.JSON(http.StatusOK, cachedResponse)
		return
	}

	characters, err := ctx.Server.Store.ListAllCharacters(ctx)
	if err != nil {
		ctx.Server.Logger.Error("failed to ListAllCharacters", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListAllCharacters : %w", err)))
		return
	}

	// レスポンスをキャッシュに保存
	if len(characters) > 0 {
		response := listAllCharactersResponse{
			Characters: characters,
		}
		// Redisへのセットが失敗しても処理を続行
		err = ctx.Server.RedisClient.Set(ctx.Context, cacheKey, response, cache.CacheDurationDay)
		if err != nil {
			ctx.Server.Logger.Warn("failed redis data set", zap.String("redis_key", cacheKey), zap.Error(err))
		}
	}

	ctx.JSON(http.StatusOK, listAllCharactersResponse{
		Characters: characters,
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
// @Router /api/v1/characters/{id} [get]
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
