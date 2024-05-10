package api

import (
	"fmt"
	"mime/multipart"
	"net/http"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/pkg/util"

	"github.com/gin-gonic/gin"
)

type listCategoriesRequest struct {
	Page int64 `form:"p"`
}

func (server *Server) ListCategories(c *gin.Context) {
	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}

	categories := []*model.Category{}

	pcates, err := server.Store.ListParentCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to server.Store.ListParentCategories : %w", err)))
		return
	}

	fmt.Println(pcates)

	// pcateごとのchild_categoriesの取得

	// pcatesをrangeで回し、categories に pcate と ccate を appned していく

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
