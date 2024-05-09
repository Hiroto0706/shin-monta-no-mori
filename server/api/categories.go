package api

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type listCategoriesRequest struct {
	Page int64 `form:"p"`
}

func (server *Server) ListCategories(c *gin.Context) {}

func (server *Server) GetCategories(c *gin.Context) {}

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
