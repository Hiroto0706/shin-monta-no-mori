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

		iccrs, err := store.ListImageChildCategoryRelationsByImageID(c, i.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
			}

			c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID() : %w", err)))
		}

		cCates := []db.ChildCategory{}
		for _, iccr := range iccrs {
			cCate, err := store.GetChildCategory(c, iccr.ChildCategoryID)
			if err != nil {
				c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetChildCategory() : %w", err)))
			}
			cCates = append(cCates, cCate)
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
func UploadImageSrc(c *gin.Context, config *util.Config, formKey string, filename string, fileType string, isSimple bool) (string, error) {
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

func DeleteImageSrc(c *gin.Context, config *util.Config, src string) error {
	storageService := &GCSStorageService{
		Config: *config,
	}
	return storageService.DeleteFile(c, src)
}

// UpdateImageCharacterRelationsIDs updates the character relations for an image.
func UpdateImageCharacterRelationsIDs(c *gin.Context, store *db.Store, imageID int64, newCharacterIDs []int64) error {
	existingRelations, err := store.ListImageCharacterRelationsByImageID(c, imageID)
	if err != nil {
		return fmt.Errorf("failed to ListImageCharacterRelationsByImageID: %w", err)
	}

	// Create a map of existing character IDs for quick lookup.
	existingIDs := make(map[int64]bool)
	for _, rel := range existingRelations {
		existingIDs[rel.CharacterID] = true
	}

	// Create a set from the new character IDs.
	newIDs := make(map[int64]bool)
	for _, id := range newCharacterIDs {
		newIDs[id] = true
	}

	// Remove relations that are not needed anymore.
	for _, rel := range existingRelations {
		if !newIDs[rel.CharacterID] {
			if err := store.DeleteImageCharacterRelations(c, rel.ID); err != nil {
				return fmt.Errorf("failed to DeleteImageCharacterRelations: %w", err)
			}
		}
	}

	// Add new relations that do not exist yet.
	for id := range newIDs {
		if !existingIDs[id] {
			_, err := store.CreateImageCharacterRelations(c, db.CreateImageCharacterRelationsParams{
				ImageID:     imageID,
				CharacterID: id,
			})
			if err != nil {
				return fmt.Errorf("failed to CreateImageCharacterRelations: %w", err)
			}
		}
	}

	return nil
}

// UpdateImageParentCategoryRelationsIDs updates the parent_category relations for an image.
func UpdateImageParentCategoryRelationsIDs(c *gin.Context, store *db.Store, imageID int64, newParentCategoryIDs []int64) error {
	existingRelations, err := store.ListImageParentCategoryRelationsByImageID(c, imageID)
	if err != nil {
		return fmt.Errorf("failed to ListImageParentCategoryRelationsByImageID: %w", err)
	}

	// Create a map of existing parent_category IDs for quick lookup.
	existingIDs := make(map[int64]bool)
	for _, rel := range existingRelations {
		existingIDs[rel.ParentCategoryID] = true
	}

	// Create a set from the new parent_category IDs.
	newIDs := make(map[int64]bool)
	for _, id := range newParentCategoryIDs {
		newIDs[id] = true
	}

	// Remove relations that are not needed anymore.
	for _, rel := range existingRelations {
		if !newIDs[rel.ParentCategoryID] {
			if err := store.DeleteImageParentCategoryRelations(c, rel.ID); err != nil {
				return fmt.Errorf("failed to DeleteImageParentCategoryRelations: %w", err)
			}
		}
	}

	// Add new relations that do not exist 	yet.
	for id := range newIDs {
		if !existingIDs[id] {
			_, err := store.CreateImageParentCategoryRelations(c, db.CreateImageParentCategoryRelationsParams{
				ImageID:          imageID,
				ParentCategoryID: id,
			})
			if err != nil {
				return fmt.Errorf("failed to CreateImageParentCategoryRelations: %w", err)
			}
		}
	}

	return nil
}

// UpdateImageChildCategoryRelationsIDs updates the child_category relations for an image.
func UpdateImageChildCategoryRelationsIDs(c *gin.Context, store *db.Store, imageID int64, newChildCategoryIDs []int64) error {
	existingRelations, err := store.ListImageChildCategoryRelationsByImageID(c, imageID)
	if err != nil {
		return fmt.Errorf("failed to ListImageChildCategoryRelationsByImageID: %w", err)
	}

	// Create a map of existing child_category IDs for quick lookup.
	existingIDs := make(map[int64]bool)
	for _, rel := range existingRelations {
		existingIDs[rel.ChildCategoryID] = true
	}

	// Create a set from the new child_category IDs.
	newIDs := make(map[int64]bool)
	for _, id := range newChildCategoryIDs {
		newIDs[id] = true
	}

	// Remove relations that are not needed anymore.
	for _, rel := range existingRelations {
		if !newIDs[rel.ChildCategoryID] {
			if err := store.DeleteImageChildCategoryRelations(c, rel.ID); err != nil {
				return fmt.Errorf("failed to DeleteImageChildCategoryRelations: %w", err)
			}
		}
	}

	// Add new relations that do not exist 	yet.
	for id := range newIDs {
		if !existingIDs[id] {
			_, err := store.CreateImageChildCategoryRelations(c, db.CreateImageChildCategoryRelationsParams{
				ImageID:         imageID,
				ChildCategoryID: id,
			})
			if err != nil {
				return fmt.Errorf("failed to CreateImageChildCategoryRelations: %w", err)
			}
		}
	}

	return nil
}
