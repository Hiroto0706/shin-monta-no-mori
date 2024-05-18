package admin

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/internal/domains/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	IMAGE_TYPE_CATEGORY = "category"
)

// TODO: 将来的にpager機能を持たせた方がいいかも？
type listCategoriesRequest struct {
	Page int64 `form:"p"`
}

// ListCategories godoc
// @Summary List categories
// @Description Retrieves a list of parent categories along with their child categories.
// @Accept  json
// @Produce  json
// @Success 200 {array} model/Category "A list of categories with parent and child category details."
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: The request is malformed or missing required fields."
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: Child categories not found for one or more parent categories."
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// @Router /api/v1/admin/categories/list [get]
func ListCategories(ctx *app.AppContext) {
	// // TODO: bind 周りの処理は関数化して共通化したほうがいい
	// var req listCategoriesRequest
	// if err := c.ShouldBindQuery(&req); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
	// 	return
	// }

	pcates, err := ctx.Server.Store.ListParentCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to ctx.Server.Store.ListParentCategories : %w", err)))
		return
	}

	categories := make([]model.Category, len(pcates))
	for i, pcate := range pcates {
		ccates, err := ctx.Server.Store.GetChildCategoriesByParentID(ctx, pcate.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID: %w", err)))
				return
			}

			ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID : %w", err)))
			return
		}

		categories[i] = model.Category{
			ParentCategory: pcate,
			ChildCategory:  ccates,
		}
	}

	ctx.JSON(http.StatusOK, categories)
}

// GetCategory godoc
// @Summary Retrieve a category
// @Description Retrieves a parent category along with its child categories by the parent category's ID
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID of the parent category to retrieve"
// @Success 200 {object} model/Category "The requested parent category with its child categories"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Failed to parse 'id' number from path parameter"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No parent category found with the given ID or no child categories found for the parent category"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to retrieve the category from the database"
// @Router /api/v1/admin/categories/{id} [get]
func GetCategory(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}

	pcate, err := ctx.Server.Store.GetParentCategory(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetParentCategory: %w", err)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetParentCategory : %w", err)))
		return
	}

	ccates, err := ctx.Server.Store.GetChildCategoriesByParentID(ctx, pcate.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID : %w", err)))
		return
	}

	category := model.NewCategory()
	category.ParentCategory = pcate
	category.ChildCategory = ccates

	ctx.JSON(http.StatusOK, category)
}

type searchCategoriesRequest struct {
	Query string `form:"q"`
}

// SearchCategories godoc
// @Summary Search parent categories
// @Description Searches for parent categories based on a query string.
// @Accept  json
// @Produce  json
// @Param   q   query   string  true  "Query string to search parent categories"
// @Success 200 {array} model/Category "List of categories with their corresponding child categories"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Failed to bind query parameters"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No child categories found for a parent category"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to retrieve categories from the database"
// @Router /api/v1/admin/categories/search [get]
func SearchCategories(ctx *app.AppContext) {
	// TODO: 親カテゴリの検索のみでなく、子カテゴリようの検索APIも追加する。
	// それか、検索機能は一つにし、親カテゴリが一致する場合は個カテゴリ全て取得、個カテゴリが一致する場合は、子カテゴリの一部と個カテゴリが持つ親カテゴリのみ取得するようにする
	var req searchCategoriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}

	q := sql.NullString{
		String: req.Query,
		Valid:  true,
	}
	pcates, err := ctx.Server.Store.SearchParentCategories(ctx, q)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to SearchParentCategories : %w", err)))
		return
	}

	categories := make([]*model.Category, len(pcates))
	for i, pcate := range pcates {
		ccates, err := ctx.Server.Store.GetChildCategoriesByParentID(ctx, pcate.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID : %w", err)))
			return
		}

		categories[i] = &model.Category{
			ParentCategory: pcate,
			ChildCategory:  ccates,
		}
	}

	ctx.JSON(http.StatusOK, categories)
}

type createParentCategoryRequest struct {
	Name      string               `form:"name" binding:"required"`
	Filename  string               `form:"filename" binding:"required"`
	ImageFile multipart.FileHeader `form:"image_file" binding:"required"`
}

// CreateParentCategory godoc
// @Summary Create a new parent category
// @Description Creates a new parent category with a name, filename, and an image file.
// @Accept  multipart/form-data
// @Produce  json
// @Param   name       formData   string  true  "Name of the parent category"
// @Param   filename   formData   string  true  "Filename for the uploaded image"
// @Param   image_file formData   file    true  "Image file for the parent category"
// @Success 200 {object} gin/H "Returns the created parent category and a success message"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or validation"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to create the parent category due to a server error"
// @Router /api/v1/admin/categories/parent/create [post]
func CreateParentCategory(ctx *app.AppContext) {
	var req createParentCategoryRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	var parentCategory db.ParentCategory
	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {
		src, err := service.UploadImageSrc(ctx.Context, &ctx.Server.Config, "image_file", req.Filename, IMAGE_TYPE_CATEGORY, false)
		if err != nil {
			return fmt.Errorf("failed to UploadImage: %w", err)
		}

		arg := db.CreateParentCategoryParams{
			Name: req.Name,
			Src:  src,
			Filename: sql.NullString{
				String: req.Filename,
				Valid:  true,
			},
		}

		parentCategory, err = ctx.Server.Store.CreateParentCategory(ctx, arg)
		if err != nil {
			return fmt.Errorf("failed to CreateParentCategory: %w", err)
		}

		return nil
	})

	if txErr != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("CreateParentCategory transaction was failed : %w", txErr)))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"parent_category": parentCategory,
		"message":         "parent_categoryの作成に成功しました",
	})
}

type editParentCategoryRequest struct {
	Name      string               `form:"name"`
	Filename  string               `form:"filename"`
	ImageFile multipart.FileHeader `form:"image_file"`
}

// EditParentCategory godoc
// @Summary Edit an existing parent category
// @Description Edits a parent category by ID, allowing updates to the category's name, filename, and associated image.
// @Accept  multipart/form-data
// @Produce  json
// @Param   id         path     int    true  "ID of the parent category to edit"
// @Param   name       formData string true  "New name of the parent category"
// @Param   filename   formData string true  "New filename for the uploaded image"
// @Param   image_file formData file   false "New image file for the parent category (optional)"
// @Success 200 {object} gin/H "Returns the updated parent category and a success message"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or missing required fields"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No parent category found with the given ID"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to update the parent category due to a server error"
// @Router /api/v1/admin/categories/parent/{id} [put]
func EditParentCategory(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}
	var req editParentCategoryRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	pcate, err := ctx.Server.Store.GetParentCategory(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetParentCategory : %w", err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetParentCategory : %w", err)))
		return
	}

	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {
		src := pcate.Src
		if pcate.Filename.String != req.Filename {
			err := service.DeleteImageSrc(ctx.Context, &ctx.Server.Config, pcate.Src)
			if err != nil {
				return err
			}

			src, err = service.UploadImageSrc(ctx.Context, &ctx.Server.Config, "image_file", req.Filename, IMAGE_TYPE_CATEGORY, false)
			if err != nil {
				return err
			}
		}

		arg := db.UpdateParentCategoryParams{
			ID:        pcate.ID,
			Name:      req.Name,
			Src:       src,
			Filename:  sql.NullString{String: pcate.Filename.String, Valid: true},
			UpdatedAt: time.Now(),
		}
		if pcate.Filename.String != req.Filename {
			arg.Filename = sql.NullString{String: req.Filename, Valid: true}
		}

		pcate, err = ctx.Server.Store.UpdateParentCategory(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("EditParentCategory transaction was failed : %w", txErr)))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"parent_category": pcate,
		"message":         "parent_categoryの編集に成功しました",
	})
}

// DeleteParentCategory godoc
// @Summary Delete a parent category
// @Description Deletes an existing parent category identified by its ID along with all its associated child categories and related image relations.
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID of the parent category to delete"
// @Success 200 {object} gin/H "Returns a success message indicating the parent category and all related entities have been deleted"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in parsing the parent category ID"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No parent category found with the given ID"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to delete the parent category or its related entities due to a server error"
// @Router /api/v1/admin/categories/parent/{id} [delete]
func DeleteParentCategory(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	pcate, err := ctx.Server.Store.GetParentCategory(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetParentCategory : %w", err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetParentCategory : %w", err)))
		return
	}

	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {
		err = service.DeleteImageSrc(ctx.Context, &ctx.Server.Config, pcate.Src)
		if err != nil {
			return fmt.Errorf("failed to DeleteImageSrc: %w", err)
		}

		// images_parent_category_relationsの削除
		err = ctx.Server.Store.DeleteAllImageParentCategoryRelationsByParentCategoryID(ctx, pcate.ID)
		if err != nil {
			return fmt.Errorf("failed to DeleteAllImageParentCategoryRelationsByParentCategoryID : %w", err)
		}

		// parent_category_idと関連するimage_child_category_relationsの削除
		ccates, err := ctx.Server.Store.GetChildCategoriesByParentID(ctx, pcate.ID)
		if err != nil {
			return fmt.Errorf("failed to GetChildCategoriesByParentID: %w", err)
		}
		for _, ccate := range ccates {
			err = ctx.Server.Store.DeleteAllImageChildCategoryRelationsByChildCategoryID(ctx, ccate.ID)
			if err != nil {
				return fmt.Errorf("failed to DeleteAllImageChildCategoryRelationsByChildCategoryID: %w", err)
			}
		}

		// 関係するchild_categoriesの全削除
		err = ctx.Server.Store.DeleteAllChildCategoriesByParentCategoryID(ctx, pcate.ID)
		if err != nil {
			return fmt.Errorf("failed to DeleteAllChildCategoriesByParentCategoryID : %w", err)
		}

		// parent_categoryの削除
		err = ctx.Server.Store.DeleteParentCategory(ctx, pcate.ID)
		if err != nil {
			return fmt.Errorf("failed to DeleteParentCategory : %w", err)
		}

		return nil
	})
	if txErr != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("DeleteParentCategory transaction was failed : %w", txErr)))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "parent_categoryの削除に成功しました",
	})
}

// CreateChildCategory godoc
// @Summary Create a new child category
// @Description Creates a new child category with a specified name and parent ID.
// @Accept  multipart/form-data
// @Produce  json
// @Param   name       formData   string  true  "Name of the child category"
// @Param   parent_id  formData   int     true  "Parent category ID to which the child category belongs"
// @Success 200 {object} gin/H "Returns the created child category along with a success message"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or missing required fields"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to create the child category due to server-side error"
// @Router /api/v1/admin/categories/child/create [post]
type createChildCategoryRequest struct {
	Name     string `form:"name" binding:"required"`
	ParentID int    `form:"parent_id" binding:"required"`
}

func CreateChildCategory(ctx *app.AppContext) {
	var req createChildCategoryRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}

	arg := db.CreateChildCategoryParams{
		Name:     req.Name,
		ParentID: int64(req.ParentID),
	}
	childCategory, err := ctx.Server.Store.CreateChildCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to CreateChildCategory : %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"child_category": childCategory,
		"message":        "child_categoryの作成に成功しました",
	})
}

type editChildCategoryRequest struct {
	Name     string `form:"name"`
	ParentID int    `form:"parent_id"`
}

// EditChildCategory godoc
// @Summary Edit a child category
// @Description Edits an existing child category identified by its ID with new name and parent ID.
// @Accept  multipart/form-data
// @Produce  json
// @Param   id        path     int    true  "ID of the child category to edit"
// @Param   name      formData string true  "New name for the child category"
// @Param   parent_id formData int    true  "New parent ID for the child category"
// @Success 200 {object} gin/H "Returns the updated child category and a success message"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in binding query parameters or the request data"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No child category found with the given ID"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to update the child category in the database"
// @Router /api/v1/admin/categories/child/{id} [put]
func EditChildCategory(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}
	var req editChildCategoryRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	ccate, err := ctx.Server.Store.GetChildCategory(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetChildCategory : %w", err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetChildCategory : %w", err)))
		return
	}

	txErr := ctx.Server.Store.ExecTx(ctx.Request.Context(), func(q *db.Queries) error {
		arg := db.UpdateChildCategoryParams{
			ID:        ccate.ID,
			Name:      req.Name,
			ParentID:  int64(req.ParentID),
			UpdatedAt: time.Now(),
		}

		ccate, err = ctx.Server.Store.UpdateChildCategory(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("EditChildCategory transaction was failed : %w", txErr)))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"child_category": ccate,
		"message":        "child_categoryの編集に成功しました",
	})
}

// DeleteChildCategory godoc
// @Summary Delete a child category
// @Description Deletes an existing child category identified by its ID.
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID of the child category to delete"
// @Success 200 {object} gin/H "Returns a success message indicating the child category has been deleted"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in parsing the child category ID"
// @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No child category found with the given ID or error in deleting the child category"
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to retrieve or delete the child category from the database"
// @Router /api/v1/admin/categories/child/{id} [delete]
func DeleteChildCategory(ctx *app.AppContext) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	ccate, err := ctx.Server.Store.GetChildCategory(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to GetChildCategory : %w", err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(fmt.Errorf("failed to GetChildCategory : %w", err)))
		return
	}

	err = ctx.Server.Store.DeleteChildCategory(ctx, ccate.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, app.ErrorResponse(fmt.Errorf("failed to DeleteChildCategory : %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "child_categoryの削除に成功しました",
	})
}
