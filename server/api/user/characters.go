package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/pkg/lib/binder"
	"strconv"

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
	var req listCharactersRequest
	if err := binder.BindQuery(ctx.Context, &req); err != nil {
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
	characters, err := ctx.Server.Store.ListAllCharacters(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListAllCharacters : %w", err)))
		return
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

		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, getCharacterResponse{
		Character: character,
	})
}
