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
	IMAGE_TYPE_IMAGE = "image"
)

func (server *Server) Greet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "new Hello World from server.",
	})
}

type listIllustrationsRequest struct {
	Page int64 `form:"p"`
}

func (server *Server) ListIllustrations(c *gin.Context) {
	var req listIllustrationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}

	illustrations := []*model.Illustration{}

	arg := db.ListImageParams{
		Limit:  int32(server.Config.ImageFetchLimit),
		Offset: int32(int(req.Page) * server.Config.ImageFetchLimit),
	}
	images, err := server.Store.ListImage(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to ListImage() : %w", err)))
			return
		}

		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to ListImage() : %w", err)))
		return
	}

	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(c, server.Store, i)

		illustrations = append(illustrations, il)
	}

	c.JSON(http.StatusOK, illustrations)
}

func (server *Server) GetIllustration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
		return
	}

	image, err := server.Store.GetImage(c, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetImage() : %w", err)))
			return
		}

		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetImage() : %w", err)))
		return
	}

	illustration := &model.Illustration{}
	illustration = service.FetchRelationInfoForIllustrations(c, server.Store, image)

	c.JSON(http.StatusOK, illustration)
}

type searchIllustrationsRequest struct {
	Page  int    `form:"p"`
	Query string `form:"q"`
}

// TODO: imageだけでなく、カテゴリでも検索ができるようにする。
// また、検索結果をtrimし、被りがないようにする
func (server *Server) SearchIllustrations(c *gin.Context) {
	var req searchIllustrationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}
	arg := db.SearchImagesParams{
		Limit:  int32(server.Config.ImageFetchLimit),
		Offset: int32(req.Page * server.Config.ImageFetchLimit),
		Title: sql.NullString{
			String: req.Query,
			Valid:  true,
		},
	}

	images, err := server.Store.SearchImages(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to SearchImages() : %w", err)))
			return
		}

		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to SearchImages() : %w", err)))
		return
	}

	illustrations := []*model.Illustration{}
	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(c, server.Store, i)

		illustrations = append(illustrations, il)
	}

	c.JSON(http.StatusOK, illustrations)
}

type createIllustrationRequest struct {
	Title             string               `json:"title" form:"title"`
	Filename          string               `json:"filename" form:"filename"`
	Characters        []int64              `json:"characters" form:"characters"`
	ParentCategories  []int64              `json:"parent_categories" form:"parent_categories"`
	ChildCategories   []int64              `json:"child_categories" form:"child_categories"`
	OriginalImageFile multipart.FileHeader `json:"original_image_file" form:"image_file"`
	SimpleImageFile   multipart.FileHeader `json:"simple_image_file" form:"image_file"`
}

func (server *Server) CreateIllustration(c *gin.Context) {
	var req createIllustrationRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	var image db.Image
	txErr := server.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
		originalSrc, err := service.UploadImage(c, &server.Config, "original_image_file", req.Filename, IMAGE_TYPE_IMAGE, false)
		if err != nil {
			return err
		}

		simpleSrc, err := service.UploadImage(c, &server.Config, "simple_image_file", req.Filename, IMAGE_TYPE_IMAGE, true)
		if err != nil {
			return err
		}

		arg := db.CreateImageParams{
			Title:       req.Title,
			OriginalSrc: originalSrc,
			SimpleSrc: sql.NullString{
				String: simpleSrc,
				Valid:  true,
			},
			OriginalFilename: req.Filename,
		}

		image, err = server.Store.CreateImage(c, arg)
		if err != nil {
			return err
		}

		// ImageCharacterRelationsの保存
		for _, c_id := range req.ChildCategories {
			arg := db.CreateImageCharacterRelationsParams{
				ImageID:     image.ID,
				CharacterID: c_id,
			}
			_, err := server.Store.CreateImageCharacterRelations(c, arg)
			if err != nil {
				return err
			}
		}

		// ImageParentCategoryRelationsの保存
		for _, pc_id := range req.ParentCategories {
			arg := db.CreateImageParentCategoryRelationsParams{
				ImageID:          image.ID,
				ParentCategoryID: pc_id,
			}

			_, err := server.Store.CreateImageParentCategoryRelations(c, arg)
			if err != nil {
				return err
			}
		}

		// ImageChildCategoryRelationsの保存
		for _, cc_id := range req.ChildCategories {
			arg := db.CreateImageChildCategoryRelationsParams{
				ImageID:         image.ID,
				ChildCategoryID: cc_id,
			}
			_, err := server.Store.CreateImageChildCategoryRelations(c, arg)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(txErr))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"title":            image.Title,
		"filename":         image.OriginalFilename,
		"characters":       req.Characters,
		"parentCategories": req.ParentCategories,
		"childCategories":  req.ChildCategories,
		"o_src":            image.OriginalSrc,
		"s_src":            image.SimpleSrc,
	})
}

type editIllustrationRequest struct {
	Title             string               `json:"title" form:"title"`
	Filename          string               `json:"filename" form:"filename"`
	Characters        []int64              `json:"characters" form:"characters"`
	ParentCategories  []int64              `json:"parent_categories" form:"parent_categories"`
	ChildCategories   []int64              `json:"child_categories" form:"child_categories"`
	OriginalImageFile multipart.FileHeader `json:"original_image_file" form:"image_file"`
	SimpleImageFile   multipart.FileHeader `json:"simple_image_file" form:"image_file"`
}

func (server *Server) EditIllustration(c *gin.Context) {}

func (server *Server) DeleteIllustration(c *gin.Context) {}
