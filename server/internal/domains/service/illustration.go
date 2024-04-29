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
		arg := sql.NullInt64{
			Int64: i.ID,
			Valid: true,
		}
		cCates, err := store.GetChildCategoriesByImageID(c, arg)
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

// isSimpleはGCSにアップロードする時に画像に'_s'をつけるために使用する
func UploadImage(c *gin.Context, config *util.Config, formKey string, filename string, fileType string, isSimple bool) (string, error) {
	f, err := c.FormFile(formKey)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", fmt.Errorf("failed to get file: %w", err)
	}

	file, err := f.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	storageService := &GCSStorageService{
		Config: *config,
	}
	return storageService.UploadFile(c, file, filename, fileType, isSimple)
}

func DeleteImage(c *gin.Context, config *util.Config, src string) error {
	storageService := &GCSStorageService{
		Config: *config,
	}
	return storageService.DeleteFile(c, src)
}
