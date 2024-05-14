package api

import (
	"database/sql"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/internal/domains/service"
	"shin-monta-no-mori/server/pkg/util"
	"strconv"
	"strings"
	"time"

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
// @Router /api/v1/admin/characters/list [get]
func (s *Server) ListCharacters(c *gin.Context) {
	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listCharactersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}

	arg := db.ListCharactersParams{
		Limit:  int32(s.Config.CharacterFetchLimit),
		Offset: int32(int(req.Page) * s.Config.CharacterFetchLimit),
	}

	characters, err := s.Store.ListCharacters(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to ListCharacters : %w", err)))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": characters,
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
func (s *Server) SearchCharacters(c *gin.Context) {
	var req searchCharactersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}
	arg := db.SearchCharactersParams{
		Limit:  int32(s.Config.CharacterFetchLimit),
		Offset: int32(req.Page * s.Config.CharacterFetchLimit),
		Query: sql.NullString{
			String: req.Query,
			Valid:  true,
		},
	}
	characters, err := s.Store.SearchCharacters(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to SearchCharacters : %w", err)))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": characters,
	})
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
func (s *Server) GetCharacter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}

	character, err := s.Store.GetCharacter(c, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetCharacter: %w", err)))
			return
		}

		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"character": character,
	})
}

type createCharacterRequest struct {
	Name      string               `form:"name" binding:"required"`
	Filename  string               `form:"filename" binding:"required"`
	ImageFile multipart.FileHeader `form:"image_file" binding:"required"`
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
func (s *Server) CreateCharacter(c *gin.Context) {
	var req createCharacterRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	var character db.Character
	txErr := s.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
		var src string
		src, err := service.UploadImageSrc(c, &s.Config, "image_file", req.Filename, IMAGE_TYPE_CHARACTER, false)
		if err != nil {
			return fmt.Errorf("failed to UploadImage: %w", err)
		}

		arg := db.CreateCharacterParams{
			Name:     req.Name,
			Src:      src,
			Filename: sql.NullString{String: req.Filename, Valid: true},
		}
		character, err = s.Store.CreateCharacter(c, arg)
		if err != nil {
			return fmt.Errorf("failed to CreateCharacter: %w", err)
		}

		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("CreateCharacter transaction was failed : %w", txErr)))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"character": character,
		"message":   "characterの作成に成功しました",
	})
}

type editCharacterRequest struct {
	Name      string               `form:"name"`
	Filename  string               `form:"filename"`
	ImageFile multipart.FileHeader `form:"image_file"`
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
func (s *Server) EditCharacter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}
	var req editCharacterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	character, err := s.Store.GetCharacter(c, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
			return
		}
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
		return
	}

	txErr := s.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
		src := character.Src
		if character.Filename.String != req.Filename {
			err := service.DeleteImageSrc(c, &s.Config, character.Src)
			if err != nil {
				return err
			}

			src, err = service.UploadImageSrc(c, &s.Config, "image_file", req.Filename, IMAGE_TYPE_CHARACTER, false)
			if err != nil {
				return err
			}
		}

		arg := db.UpdateCharacterParams{
			ID:        character.ID,
			Name:      req.Name,
			Src:       src,
			Filename:  sql.NullString{String: character.Filename.String, Valid: true},
			UpdatedAt: time.Now(),
		}
		if character.Filename.String != req.Filename {
			arg.Filename = sql.NullString{String: req.Filename, Valid: true}
		}

		character, err = s.Store.UpdateCharacter(c, arg)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("EditCharacter transaction was failed : %w", txErr)))
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
func (s *Server) DeleteCharacter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}
	character, err := s.Store.GetCharacter(c, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
			return
		}
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetCharacter : %w", err)))
		return
	}

	txErr := s.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
		err = service.DeleteImageSrc(c, &s.Config, character.Src)
		if err != nil {
			return fmt.Errorf("failed to DeleteImageSrc: %w", err)
		}

		// images_character_relationsの削除
		err = s.Store.DeleteAllImageCharacterRelationsByCharacterID(c, character.ID)
		if err != nil {
			return fmt.Errorf("failed to DeleteAllImageCharacterRelationsByCharacterID : %w", err)
		}

		err = s.Store.DeleteCharacter(c, int64(id))
		if err != nil {
			return fmt.Errorf("failed to DeleteCharacter: %w", err)
		}
		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("DeleteCharacter transaction was failed : %w", txErr)))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "characterの削除に成功しました",
	})
}
