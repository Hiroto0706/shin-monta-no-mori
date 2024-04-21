package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) Greet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "new Hello World from server.",
	})
}

func (server *Server) ListIllustrations(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("p"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Errorf("failed to parse page number from query param : %w", err)))
	}

	illustrations := []*model.Illustration{}

	arg := db.ListImageParams{
		Limit:  int32(server.Config.ImageFetchLimit),
		Offset: int32(page * server.Config.ImageFetchLimit),
	}
	images, err := server.Store.ListImage(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, NewErrorResponse(fmt.Errorf("failed to ListImage() : %w", err)))
		}

		c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("failed to ListImage() : %w", err)))
	}

	for _, i := range images {
		// キャラクターの取得
		icrs, err := server.Store.ListImageCharacterRelationsByImageID(c, i.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
			}

			c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
		}

		log.Println(icrs)
		characters := []db.Character{}
		for _, icr := range icrs {
			char, err := server.Store.GetCharacter(c, icr.CharacterID)
			if err != nil {
				c.JSON(http.StatusNotFound, NewErrorResponse(fmt.Errorf("failed to GetCharacter() : %w", err)))
			}

			characters = append(characters, char)
		}

		// カテゴリーの取得
		_, err = server.Store.ListImageParentCategoryRelationsByImageID(c, i.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
			}

			c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
		}

		categories := []*model.Category{}
		// for _, ipcr := range ipcrs {
		// 	pCate, err := server.Store.GetParentCategory(c, ipcr.ParentCategoryID)
		// 	if err != nil {
		// 		c.JSON(http.StatusNotFound, NewErrorResponse(fmt.Errorf("failed to GetParentCategory() : %w", err)))
		// 	}
		// 	cCates , err := server.Store.

		// pc := model.NewCategory()
		// pc.ParentCategory = pCate

		// categories = append(categories, pc)
		// }

		il := model.NewIllustration()

		il.Image = i
		il.Character = characters
		il.Category = categories

		illustrations = append(illustrations, il)
	}

	c.JSON(http.StatusOK, gin.H{
		"illustrations": illustrations,
	})
}
