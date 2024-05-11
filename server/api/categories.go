package api

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/internal/domains/service"
	"shin-monta-no-mori/server/pkg/util"
	"strconv"
	"strings"

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
func (server *Server) ListCategories(c *gin.Context) {
	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}

	pcates, err := server.Store.ListParentCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to server.Store.ListParentCategories : %w", err)))
		return
	}

	categories := make([]*model.Category, len(pcates))
	for i, pcate := range pcates {
		ccates, err := server.Store.GetChildCategoriesByParentID(c, pcate.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID: %w", err)))
				return
			}

			c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID : %w", err)))
			return
		}

		categories[i] = &model.Category{
			ParentCategory: pcate,
			ChildCategory:  ccates,
		}
	}

	c.JSON(http.StatusOK, categories)
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
func (server *Server) GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}

	pcate, err := server.Store.GetParentCategory(c, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetParentCategory: %w", err)))
			return
		}

		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetParentCategory : %w", err)))
		return
	}

	ccates, err := server.Store.GetChildCategoriesByParentID(c, pcate.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID : %w", err)))
		return
	}

	category := model.NewCategory()
	category.ParentCategory = pcate
	category.ChildCategory = ccates

	c.JSON(http.StatusOK, category)
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
func (server *Server) SearchCategories(c *gin.Context) {
	var req searchCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return
	}

	q := sql.NullString{
		String: req.Query,
		Valid:  true,
	}
	pcates, err := server.Store.SearchParentCategories(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to SearchParentCategories : %w", err)))
		return
	}

	categories := make([]*model.Category, len(pcates))
	for i, pcate := range pcates {
		ccates, err := server.Store.GetChildCategoriesByParentID(c, pcate.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID: %w", err)))
				return
			}

			c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetChildCategoriesByParentID : %w", err)))
			return
		}

		categories[i] = &model.Category{
			ParentCategory: pcate,
			ChildCategory:  ccates,
		}
	}

	c.JSON(http.StatusOK, categories)
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
func (server *Server) CreateParentCategory(c *gin.Context) {
	var req createParentCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	var parentCategory db.ParentCategory
	txErr := server.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
		src, err := service.UploadImageSrc(c, &server.Config, "image_file", req.Filename, IMAGE_TYPE_CATEGORY, false)
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

		parentCategory, err = server.Store.CreateParentCategory(c, arg)
		if err != nil {
			return fmt.Errorf("failed to CreateParentCategory: %w", err)
		}

		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("CreateParentCategory transaction was failed : %w", txErr)))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"parent_category": parentCategory,
		"message":         "categoryの作成に成功しました",
	})
}

type editCategoryRequest struct {
	Name      string               `form:"name"`
	Filename  string               `form:"filename"`
	ImageFile multipart.FileHeader `form:"image_file"`
}

func (server *Server) EditCategory(c *gin.Context) {}

func (server *Server) DeleteCategory(c *gin.Context) {}
