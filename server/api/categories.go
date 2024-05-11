package api

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/pkg/util"

	"github.com/gin-gonic/gin"
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

func (server *Server) GetCategory(c *gin.Context) {}

type searchCategoriesRequest struct {
	Page  int    `form:"p"`
	Query string `form:"q"`
}

func (server *Server) SearchCategories(c *gin.Context) {}

type createCategoryRequest struct {
	Title     string               `form:"title" binding:"required"`
	Filename  string               `form:"filename" binding:"required"`
	ImageFile multipart.FileHeader `form:"image_file" binding:"required"`
}

func (server *Server) CreateCategory(c *gin.Context) {}

type editCategoryRequest struct {
	Title     string               `form:"title"`
	Filename  string               `form:"filename"`
	ImageFile multipart.FileHeader `form:"image_file"`
}

func (server *Server) EditCategory(c *gin.Context) {}

func (server *Server) DeleteCategory(c *gin.Context) {}
