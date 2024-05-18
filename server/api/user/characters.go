package user

import (
	"fmt"
	"net/http"
	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"

	"github.com/gin-gonic/gin"
)

const (
	IMAGE_TYPE_CHARACTER = "character"
)

type listCharactersRequest struct {
	Page int64 `form:"p"`
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
// @Router /api/v1/characters/list [get]
func ListCharacters(ctx *app.AppContext) {
	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listCharactersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}

	arg := db.ListCharactersParams{
		Limit:  int32(ctx.Server.Config.CharacterFetchLimit),
		Offset: int32(int(req.Page) * ctx.Server.Config.CharacterFetchLimit),
	}

	characters, err := ctx.Server.Store.ListCharacters(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListCharacters : %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"characters": characters,
	})
}
