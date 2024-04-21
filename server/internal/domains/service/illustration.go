package service

import (
	"database/sql"
	"fmt"
	"net/http"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/pkg/util"

	"github.com/gin-gonic/gin"
)

func FetchRelationInfoForIllustrations(c *gin.Context, store *db.Store, i db.Image) *model.Illustration {
	// キャラクターの取得
	icrs, err := store.ListImageCharacterRelationsByImageID(c, i.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
		}

		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
	}

	characters := []db.Character{}
	for _, icr := range icrs {
		char, err := store.GetCharacter(c, icr.CharacterID)
		if err != nil {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetCharacter() : %w", err)))
		}

		characters = append(characters, char)
	}

	// カテゴリーの取得
	ipcrs, err := store.ListImageParentCategoryRelationsByImageID(c, i.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
		}

		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
	}

	categories := []*model.Category{}
	for _, ipcr := range ipcrs {
		pCate, err := store.GetParentCategory(c, ipcr.ParentCategoryID)
		if err != nil {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetParentCategory() : %w", err)))
		}
		cCates, err := store.GetChildCategoriesByParentID(c, ipcr.ParentCategoryID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
			}

			c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
		}

		cate := model.NewCategory()
		cate.ParentCategory = pCate
		cate.ChildCategory = cCates

		categories = append(categories, cate)
	}

	il := model.NewIllustration()

	il.Image = i
	il.Character = characters
	il.Category = categories

	return il
}
