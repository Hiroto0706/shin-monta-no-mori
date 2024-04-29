package api

import (
	"database/sql"
	"fmt"
	"log"
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
	Title             string               `json:"title" form:"title" binding:"required"`
	Filename          string               `json:"filename" form:"filename" binding:"required"`
	Characters        []int64              `form:"characters[]"`
	ParentCategories  []int64              `form:"parent_categories[]"`
	ChildCategories   []int64              `form:"child_categories[]"`
	OriginalImageFile multipart.FileHeader `json:"original_image_file" form:"image_file" binding:"required"`
	SimpleImageFile   multipart.FileHeader `json:"simple_image_file" form:"image_file"`
}

func (server *Server) CreateIllustration(c *gin.Context) {
	var req createIllustrationRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}
	log.Println("req", req)
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	var image db.Image
	txErr := server.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
		originalSrc, err := service.UploadImageSrc(c, &server.Config, "original_image_file", req.Filename, IMAGE_TYPE_IMAGE, false)
		if err != nil {
			return fmt.Errorf("failed to UploadImage: %w", err)
		}

		simpleSrc, err := service.UploadImageSrc(c, &server.Config, "simple_image_file", req.Filename, IMAGE_TYPE_IMAGE, true)
		if err != nil {
			return fmt.Errorf("failed to UploadImage: %w", err)
		}

		arg := db.CreateImageParams{
			Title:            req.Title,
			OriginalSrc:      originalSrc,
			OriginalFilename: req.Filename,
			SimpleSrc:        sql.NullString{String: "", Valid: false},
			SimpleFilename:   sql.NullString{String: "", Valid: false},
		}
		if simpleSrc != "" {
			arg.SimpleSrc = sql.NullString{String: simpleSrc, Valid: true}
			arg.SimpleFilename = sql.NullString{String: req.Filename + "_s", Valid: true}
		}

		image, err = server.Store.CreateImage(c, arg)
		if err != nil {
			return fmt.Errorf("failed to CreateImage: %w", err)
		}

		// ImageCharacterRelationsの保存
		for _, c_id := range req.Characters {
			arg := db.CreateImageCharacterRelationsParams{
				ImageID:     image.ID,
				CharacterID: c_id,
			}
			_, err := server.Store.CreateImageCharacterRelations(c, arg)
			if err != nil {
				return fmt.Errorf("failed to CreateImageCharacterRelations: %w", err)
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
				return fmt.Errorf("failed to CreateImageParentCategoryRelations: %w", err)
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
				return fmt.Errorf("failed to CreateImageChildCategoryRelations: %w", err)
			}
		}

		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("CreateImage transaction was failed : %w", txErr)))
		return
	}

	image, err := server.Store.GetImage(c, int64(image.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetImage() : %w", err)))
			return
		}

		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetImage() : %w", err)))
		return
	}

	illustration := service.FetchRelationInfoForIllustrations(c, server.Store, image)

	c.JSON(http.StatusOK, gin.H{
		"title":      illustration.Image.Title,
		"filename":   illustration.Image.OriginalFilename,
		"characters": illustration.Character,
		"category":   illustration.Category,
		"o_src":      illustration.Image.OriginalSrc,
		"s_src":      illustration.Image.SimpleSrc,
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

func (server *Server) EditIllustration(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
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

	illustration := service.FetchRelationInfoForIllustrations(c, server.Store, image)

	var req editIllustrationRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}
	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

	txErr := server.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
		// TODO: ファイル名だけ変わるorファイル名そのままで画像だけ変わる場合は既存の画像を削除し、更新するようにする
		originalSrc, err := service.UploadImageSrc(c, &server.Config, "original_image_file", req.Filename, IMAGE_TYPE_IMAGE, false)
		if err != nil {
			return err
		}

		// TODO: ファイル名だけ変わるorファイル名そのままで画像だけ変わる場合は既存の画像を削除し、更新するようにする
		simpleSrc, err := service.UploadImageSrc(c, &server.Config, "simple_image_file", req.Filename, IMAGE_TYPE_IMAGE, true)
		if err != nil {
			return err
		}

		// TODO: UPDATEに変える
		arg := db.CreateImageParams{
			Title:       req.Title,
			OriginalSrc: originalSrc,
			SimpleSrc: sql.NullString{
				String: simpleSrc,
				Valid:  true,
			},
			OriginalFilename: req.Filename,
		}

		// TODO: UPDATEに変える
		image, err = server.Store.CreateImage(c, arg)
		if err != nil {
			return err
		}

		// TODO: UPDATEに変える
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

		// TODO: UPDATEに変える
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

		// TODO: UPDATEに変える
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

	// image, err = server.Store.GetImage(c, int64(image.ID))
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetImage() : %w", err)))
	// 		return
	// 	}

	// 	c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetImage() : %w", err)))
	// 	return
	// }

	// illustration = service.FetchRelationInfoForIllustrations(c, server.Store, image)

	c.JSON(http.StatusOK, gin.H{
		"title":      illustration.Image.Title,
		"filename":   illustration.Image.OriginalFilename,
		"characters": illustration.Character,
		"category":   illustration.Category,
		"o_src":      illustration.Image.OriginalSrc,
		"s_src":      illustration.Image.SimpleSrc,
	})
}

func (server *Server) DeleteIllustration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
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

	txErr := server.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
		err = service.DeleteImageSrc(c, &server.Config, image.OriginalSrc)
		if err != nil {
			return fmt.Errorf("failed to DeleteImageSrc: %w", err)
		}

		err = service.DeleteImageSrc(c, &server.Config, image.SimpleSrc.String)
		if err != nil {
			return fmt.Errorf("failed to DeleteImageSrc: %w", err)
		}

		// image_child_category_relationsを削除
		iccrs, err := server.Store.ListImageChildCategoryRelationsByImageID(c, image.ID)
		if err != nil {
			return fmt.Errorf("failed to ListImageChildCategoryRelationsByImageID: %w", err)
		}
		for _, iccr := range iccrs {
			err = server.Store.DeleteImageChildCategoryRelations(c, iccr.ID)
			if err != nil {
				return fmt.Errorf("failed to DeleteImageChildCategoryRelations: %w", err)
			}
		}

		// image_parent_category_relationsを削除
		ipcrs, err := server.Store.ListImageParentCategoryRelationsByImageID(c, image.ID)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("failed to ListImageParentCategoryRelationsByImageID: %w", err)
		}
		for _, ipcr := range ipcrs {
			err = server.Store.DeleteImageParentCategoryRelations(c, ipcr.ID)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("failed to DeleteImageParentCategoryRelations: %w", err)
			}
		}

		// image_character_relationsを削除
		icrs, err := server.Store.ListImageCharacterRelationsByImageID(c, image.ID)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("failed to ListImageCharacterRelationsByImageID: %w", err)
		}
		for _, icr := range icrs {
			err = server.Store.DeleteImageCharacterRelations(c, icr.ID)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("failed to DeleteImageCharacterRelations: %w", err)
			}
		}

		// Imageを削除
		err = server.Store.DeleteImage(c, image.ID)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("failed to DeleteImage: %w", err)
		}

		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(txErr))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted image",
	})
}
