package service

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"shin-monta-no-mori/internal/app"
	db "shin-monta-no-mori/internal/db/sqlc"
	model "shin-monta-no-mori/internal/domains/models"
	"shin-monta-no-mori/pkg/util"

	"github.com/gin-gonic/gin"
)

const (
	pngExtension = ".png"
)

// FetchRelationInfoForIllustrations はimageと関連するcharacterやcategoryを取得する処理
func FetchRelationInfoForIllustrations(c *gin.Context, store *db.Store, i db.Image) *model.Illustration {
	// キャラクターの取得
	icrs, err := store.ListImageCharacterRelationsByImageID(c, i.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID : %w", err)))
		}

		c.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID : %w", err)))
	}

	characters := []*model.Character{}
	for _, icr := range icrs {
		char, err := store.GetCharacter(c, icr.CharacterID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetImage: %w", err)))
			}

			c.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
		}
		character := &model.Character{Character: char}

		characters = append(characters, character)
	}

	// image.IDに関連するparent_categoryの取得
	ipcrs, err := store.ListImageParentCategoryRelationsByImageID(c, i.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID : %w", err)))
	}
	pCates := []db.ParentCategory{}
	for _, ipcr := range ipcrs {
		pCate, err := store.GetParentCategory(c, ipcr.ParentCategoryID)
		if err != nil {
			c.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetParentCategory : %w", err)))
		}
		pCates = append(pCates, pCate)
	}

	// image.IDに関連するchild_categoryの取得
	iccrs, err := store.ListImageChildCategoryRelationsByImageID(c, i.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ListImageCharacterRelationsByImageID : %w", err)))
	}
	cCates := []db.ChildCategory{}
	for _, iccr := range iccrs {
		cCate, err := store.GetChildCategory(c, iccr.ChildCategoryID)
		if err != nil {
			c.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetChildCategory : %w", err)))
		}
		cCates = append(cCates, cCate)
	}

	categories := []*model.Category{}
	for _, pCate := range pCates {
		cate := model.NewCategory()
		cate.ParentCategory = pCate
		for _, cCate := range cCates {
			if cCate.ParentID == pCate.ID {
				cate.ChildCategory = append(cate.ChildCategory, cCate)
			}
		}
		categories = append(categories, cate)
	}

	il := model.NewIllustration()

	il.Image = i
	il.Characters = characters
	il.Categories = categories

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

	ext := strings.ToLower(filepath.Ext(f.Filename))
	if ext != pngExtension {
		return "", errors.New("please upload only png extension image")
	}

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
func UpdateImageCharacterRelationsIDs(c *gin.Context, store *db.Store, imageID int64, requestCharacterIDs []int64) error {
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
	requestIDs := make(map[int64]bool)
	for _, id := range requestCharacterIDs {
		requestIDs[id] = true
	}

	// Remove relations that are not needed anymore.
	for _, existRel := range existingRelations {
		if !requestIDs[existRel.CharacterID] {
			if err := store.DeleteImageCharacterRelations(c, existRel.ID); err != nil {
				return fmt.Errorf("failed to DeleteImageCharacterRelations: %w", err)
			}
		}
	}

	// Add new relations that do not exist yet.
	for requestID := range requestIDs {
		if !existingIDs[requestID] {
			_, err := store.CreateImageCharacterRelations(c, db.CreateImageCharacterRelationsParams{
				ImageID:     imageID,
				CharacterID: requestID,
			})
			if err != nil {
				return fmt.Errorf("failed to CreateImageCharacterRelations: %w", err)
			}
		}
	}

	return nil
}

// UpdateImageParentCategoryRelationsIDs updates the parent_category relations for an image.
func UpdateImageParentCategoryRelationsIDs(c *gin.Context, store *db.Store, imageID int64, requestParentCategoryIDs []int64) error {
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
	requestIDs := make(map[int64]bool)
	for _, id := range requestParentCategoryIDs {
		requestIDs[id] = true
	}

	// Remove relations that are not needed anymore.
	for _, existRel := range existingRelations {
		if !requestIDs[existRel.ParentCategoryID] {
			if err := store.DeleteImageParentCategoryRelations(c, existRel.ID); err != nil {
				return fmt.Errorf("failed to DeleteImageParentCategoryRelations: %w", err)
			}
		}
	}

	// Add new relations that do not exist 	yet.
	for requestID := range requestIDs {
		if !existingIDs[requestID] {
			_, err := store.CreateImageParentCategoryRelations(c, db.CreateImageParentCategoryRelationsParams{
				ImageID:          imageID,
				ParentCategoryID: requestID,
			})
			if err != nil {
				return fmt.Errorf("failed to CreateImageParentCategoryRelations: %w", err)
			}
		}
	}

	return nil
}

// UpdateImageChildCategoryRelationsIDs updates the child_category relations for an image.
func UpdateImageChildCategoryRelationsIDs(c *gin.Context, store *db.Store, imageID int64, requestChildCategoryIDs []int64) error {
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
	requestIDs := make(map[int64]bool)
	for _, id := range requestChildCategoryIDs {
		requestIDs[id] = true
	}

	// Remove relations that are not needed anymore.
	for _, existRel := range existingRelations {
		if !requestIDs[existRel.ChildCategoryID] {
			if err := store.DeleteImageChildCategoryRelations(c, existRel.ID); err != nil {
				return fmt.Errorf("failed to DeleteImageChildCategoryRelations: %w", err)
			}
		}
	}

	// Add new relations that do not exist 	yet.
	for requestID := range requestIDs {
		if !existingIDs[requestID] {
			_, err := store.CreateImageChildCategoryRelations(c, db.CreateImageChildCategoryRelationsParams{
				ImageID:         imageID,
				ChildCategoryID: requestID,
			})
			if err != nil {
				return fmt.Errorf("failed to CreateImageChildCategoryRelations: %w", err)
			}
		}
	}

	return nil
}
