package api

// import (
// 	"database/sql"
// 	"fmt"
// 	"mime/multipart"
// 	"net/http"
// 	db "shin-monta-no-mori/server/internal/db/sqlc"
// 	model "shin-monta-no-mori/server/internal/domains/models"
// 	"shin-monta-no-mori/server/internal/domains/service"
// 	"shin-monta-no-mori/server/pkg/util"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// const (
// 	IMAGE_TYPE_IMAGE = "image"
// )

// func (s *Server) Greet(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "new Hello World from server.",
// 	})
// }

// type listIllustrationsRequest struct {
// 	Page int64 `form:"p"`
// }

// // ListIllustrations godoc
// // @Summary List illustrations
// // @Description Retrieves a paginated list of illustrations based on the provided page number.
// // @Accept  json
// // @Produce  json
// // @Param   p   query   int  true  "Page number for pagination"
// // @Success 200 {array} model/Illustration "A list of illustrations"
// // @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: The request is malformed or missing required fields."
// // @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// // @Router /api/v1/admin/illustrations/list [get]
// func (s *Server) ListIllustrations(c *gin.Context) {
// 	// TODO: bind 周りの処理は関数化して共通化したほうがいい
// 	var req listIllustrationsRequest
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
// 		return
// 	}

// 	illustrations := []*model.Illustration{}

// 	arg := db.ListImageParams{
// 		Limit:  int32(s.Config.ImageFetchLimit),
// 		Offset: int32(int(req.Page) * s.Config.ImageFetchLimit),
// 	}
// 	images, err := s.Store.ListImage(c, arg)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to ListImage : %w", err)))
// 		return
// 	}

// 	for _, i := range images {
// 		il := service.FetchRelationInfoForIllustrations(c, s.Store, i)

// 		illustrations = append(illustrations, il)
// 	}

// 	c.JSON(http.StatusOK, illustrations)
// }

// // GetIllustration godoc
// // @Summary Retrieve an illustration
// // @Description Retrieves a single illustration by its ID
// // @Accept  json
// // @Produce  json
// // @Param   id   path   int  true  "ID of the illustration to retrieve"
// // @Success 200 {object} model/Illustration "The requested illustration"
// // @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Failed to parse 'id' number from path parameter"
// // @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No illustration found with the given ID"
// // @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to retrieve the illustration from the database"
// // @Router /api/v1/admin/illustrations/{id} [get]
// func (s *Server) GetIllustration(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, util.NewErrorResponse(fmt.Errorf("failed to parse 'id' number from from path parameter : %w", err)))
// 		return
// 	}

// 	image, err := s.Store.GetImage(c, int64(id))
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetImage: %w", err)))
// 			return
// 		}

// 		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
// 		return
// 	}

// 	illustration := &model.Illustration{}
// 	illustration = service.FetchRelationInfoForIllustrations(c, s.Store, image)

// 	c.JSON(http.StatusOK, illustration)
// }

// type searchIllustrationsRequest struct {
// 	Page  int    `form:"p"`
// 	Query string `form:"q"`
// }

// // TODO: imageだけでなく、カテゴリでも検索ができるようにする。
// // また、検索結果をtrimし、被りがないようにする
// // SearchIllustrations godoc
// // @Summary Search illustrations
// // @Description Searches for illustrations based on a query and page number.
// // @Accept  json
// // @Produce  json
// // @Param   p     query   int    true  "Page number for pagination"
// // @Param   q     query   string true  "Query string for searching illustrations by title or category"
// // @Success 200   {array} model/Illustration "List of matched illustrations"
// // @Failure 400   {object} request/JSONResponse{data=string} "Bad Request: The request is malformed or missing required fields."
// // @Failure 500   {object} request/JSONResponse{data=string} "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// // @Router /api/v1/admin/illustrations/search [get]
// func (s *Server) SearchIllustrations(c *gin.Context) {
// 	var req searchIllustrationsRequest
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
// 		return
// 	}
// 	arg := db.SearchImagesParams{
// 		Limit:  int32(s.Config.ImageFetchLimit),
// 		Offset: int32(req.Page * s.Config.ImageFetchLimit),
// 		Query: sql.NullString{
// 			String: req.Query,
// 			Valid:  true,
// 		},
// 	}

// 	images, err := s.Store.SearchImages(c, arg)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to SearchImages : %w", err)))
// 		return
// 	}

// 	illustrations := []*model.Illustration{}
// 	for _, i := range images {
// 		il := service.FetchRelationInfoForIllustrations(c, s.Store, i)

// 		illustrations = append(illustrations, il)
// 	}

// 	c.JSON(http.StatusOK, illustrations)
// }

// type createIllustrationRequest struct {
// 	Title             string               `form:"title" binding:"required"`
// 	Filename          string               `form:"filename" binding:"required"`
// 	Characters        []int64              `form:"characters[]"`
// 	ParentCategories  []int64              `form:"parent_categories[]"`
// 	ChildCategories   []int64              `form:"child_categories[]"`
// 	OriginalImageFile multipart.FileHeader `form:"original_image_file" binding:"required"`
// 	SimpleImageFile   multipart.FileHeader `form:"simple_image_file"`
// }

// // CreateIllustration godoc
// // @Summary Create a new illustration
// // @Description Creates a new illustration with title, filename, characters, categories, and image files.
// // @Tags illustrations
// // @Accept  multipart/form-data
// // @Produce  json
// // @Param   title              formData   string                 true  "Title of the illustration"
// // @Param   filename           formData   string                 true  "Filename for the uploaded image"
// // @Param   characters[]       formData   []int64                true  "List of character IDs associated with the illustration"
// // @Param   parent_categories[] formData []int64                 true  "List of parent category IDs associated with the illustration"
// // @Param   child_categories[] formData   []int64                true  "List of child category IDs associated with the illustration"
// // @Param   original_image_file formData file                   true  "Original image file for the illustration"
// // @Param   simple_image_file  formData  file                   false "Simple image file for the illustration (optional)"
// // @Success 200 {object} gin/H "Returns the created illustration and a success message"
// // @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or validation"
// // @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to create the illustration due to a server error"
// // @Router /api/v1/admin/illustrations/create [post]
// func (s *Server) CreateIllustration(c *gin.Context) {
// 	var req createIllustrationRequest
// 	if err := c.ShouldBind(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
// 		return
// 	}
// 	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

// 	var image db.Image
// 	txErr := s.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {

// 		var err error
// 		var originalSrc string
// 		if req.Filename != "" {
// 			originalSrc, err = service.UploadImageSrc(c, &s.Config, "original_image_file", req.Filename, IMAGE_TYPE_IMAGE, false)
// 			if err != nil {
// 				return fmt.Errorf("failed to UploadImage: %w", err)
// 			}
// 		}

// 		var simpleSrc string
// 		if req.Filename != "" && req.SimpleImageFile.Size != 0 {
// 			simpleSrc, err = service.UploadImageSrc(c, &s.Config, "simple_image_file", req.Filename, IMAGE_TYPE_IMAGE, true)
// 			if err != nil {
// 				return fmt.Errorf("failed to UploadImage: %w", err)
// 			}
// 		}

// 		arg := db.CreateImageParams{
// 			Title:            req.Title,
// 			OriginalSrc:      originalSrc,
// 			OriginalFilename: req.Filename,
// 			SimpleSrc:        sql.NullString{String: "", Valid: false},
// 			SimpleFilename:   sql.NullString{String: "", Valid: false},
// 		}
// 		if simpleSrc != "" {
// 			arg.SimpleSrc = sql.NullString{String: simpleSrc, Valid: true}
// 			arg.SimpleFilename = sql.NullString{String: req.Filename + "_s", Valid: true}
// 		}

// 		image, err = s.Store.CreateImage(c, arg)
// 		if err != nil {
// 			return fmt.Errorf("failed to CreateImage: %w", err)
// 		}

// 		// ImageCharacterRelationsの保存
// 		for _, c_id := range req.Characters {
// 			arg := db.CreateImageCharacterRelationsParams{
// 				ImageID:     image.ID,
// 				CharacterID: c_id,
// 			}
// 			_, err := s.Store.CreateImageCharacterRelations(c, arg)
// 			if err != nil {
// 				return fmt.Errorf("failed to CreateImageCharacterRelations: %w", err)
// 			}
// 		}

// 		// ImageParentCategoryRelationsの保存
// 		for _, pc_id := range req.ParentCategories {
// 			arg := db.CreateImageParentCategoryRelationsParams{
// 				ImageID:          image.ID,
// 				ParentCategoryID: pc_id,
// 			}

// 			_, err := s.Store.CreateImageParentCategoryRelations(c, arg)
// 			if err != nil {
// 				return fmt.Errorf("failed to CreateImageParentCategoryRelations: %w", err)
// 			}
// 		}

// 		// ImageChildCategoryRelationsの保存
// 		for _, cc_id := range req.ChildCategories {
// 			arg := db.CreateImageChildCategoryRelationsParams{
// 				ImageID:         image.ID,
// 				ChildCategoryID: cc_id,
// 			}
// 			_, err := s.Store.CreateImageChildCategoryRelations(c, arg)
// 			if err != nil {
// 				return fmt.Errorf("failed to CreateImageChildCategoryRelations: %w", err)
// 			}
// 		}

// 		return nil
// 	})

// 	if txErr != nil {
// 		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("CreateImage transaction was failed : %w", txErr)))
// 		return
// 	}

// 	image, err := s.Store.GetImage(c, int64(image.ID))
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetImage() : %w", err)))
// 			return
// 		}

// 		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetImage() : %w", err)))
// 		return
// 	}

// 	illustration := service.FetchRelationInfoForIllustrations(c, s.Store, image)

// 	c.JSON(http.StatusOK, gin.H{
// 		"illustration": illustration,
// 		"message":      "illustrationの作成に成功しました",
// 	})
// }

// type editIllustrationRequest struct {
// 	Title             string               `form:"title"`
// 	Filename          string               `form:"filename"`
// 	Characters        []int64              `form:"characters[]"`
// 	ParentCategories  []int64              `form:"parent_categories[]"`
// 	ChildCategories   []int64              `form:"child_categories[]"`
// 	OriginalImageFile multipart.FileHeader `form:"original_image_file"`
// 	SimpleImageFile   multipart.FileHeader `form:"simple_image_file"`
// }

// // EditIllustration godoc
// // @Summary Edit an illustration
// // @Description Updates an illustration by its ID with new title, filename, and optionally updates the image.
// // @Accept  multipart/form-data
// // @Produce  json
// // @Param   id          path     int    true  "ID of the illustration to update"
// // @Param   title       formData string true  "New title of the illustration"
// // @Param   filename    formData string true  "New filename for the illustration; used in image re-upload"
// // @Param   image_file  formData file   false "New image file for the illustration"
// // @Param   characters  formData []int  false "List of character IDs associated with the illustration"
// // @Param   parentCategories formData []int false "List of parent category IDs associated with the illustration"
// // @Param   childCategories  formData []int false "List of child category IDs associated with the illustration"
// // @Success 200 {object} gin/H "Returns the updated illustration and a success message"
// // @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error in data binding or validation"
// // @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No illustration found with the given ID"
// // @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to update the illustration due to a server error"
// // @Router /api/v1/admin/illustrations/{id} [put]
// func (s *Server) EditIllustration(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
// 		return
// 	}
// 	var req editIllustrationRequest
// 	if err := c.ShouldBind(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
// 		return
// 	}
// 	req.Filename = strings.ReplaceAll(req.Filename, " ", "-")

// 	image, err := s.Store.GetImage(c, int64(id))
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
// 		return
// 	}

// 	txErr := s.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
// 		originalSrc := image.OriginalSrc
// 		if image.OriginalFilename != req.Filename {
// 			err := service.DeleteImageSrc(c, &s.Config, image.OriginalSrc)
// 			if err != nil {
// 				return err
// 			}

// 			originalSrc, err = service.UploadImageSrc(c, &s.Config, "original_image_file", req.Filename, IMAGE_TYPE_IMAGE, false)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		simpleSrc := image.SimpleSrc.String
// 		if image.OriginalFilename != req.Filename && image.SimpleFilename.String != "" {
// 			err := service.DeleteImageSrc(c, &s.Config, image.SimpleSrc.String)
// 			if err != nil {
// 				return err
// 			}

// 			simpleSrc, err = service.UploadImageSrc(c, &s.Config, "simple_image_file", req.Filename, IMAGE_TYPE_IMAGE, true)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// imageのUpdate処理
// 		arg := db.UpdateImageParams{
// 			ID:               image.ID,
// 			Title:            req.Title,
// 			OriginalSrc:      originalSrc,
// 			SimpleSrc:        sql.NullString{String: "", Valid: false},
// 			OriginalFilename: req.Filename,
// 			SimpleFilename:   sql.NullString{String: "", Valid: false},
// 			// TODO: timezoneがUTCになっている。厳密な時系列を扱う必要がある課題が出た時に修正する必要あり。
// 			UpdatedAt: time.Now(),
// 		}
// 		if simpleSrc != "" {
// 			arg.SimpleSrc = sql.NullString{String: simpleSrc, Valid: true}
// 			arg.SimpleFilename = sql.NullString{String: req.Filename + "_s", Valid: true}
// 		}
// 		image, err = s.Store.UpdateImage(c, arg)
// 		if err != nil {
// 			return err
// 		}

// 		// TODO: relation周りのUpdate処理は共通化できそう
// 		// image_character_relationsのUpdate処理
// 		err = service.UpdateImageCharacterRelationsIDs(c, s.Store, image.ID, req.Characters)
// 		if err != nil {
// 			return fmt.Errorf("failed to server.UpdateImageCharacterRelationsIDs: %w", err)
// 		}

// 		// image_parent_category_relationsのUpdate処理
// 		err = service.UpdateImageParentCategoryRelationsIDs(c, s.Store, image.ID, req.ParentCategories)
// 		if err != nil {
// 			return fmt.Errorf("failed to server.UpdateImageParentCategoryRelationsIDs: %w", err)
// 		}

// 		// image_child_category_relationsのUpdate処理
// 		err = service.UpdateImageChildCategoryRelationsIDs(c, s.Store, image.ID, req.ChildCategories)
// 		if err != nil {
// 			return fmt.Errorf("failed to server.UpdateImageChildCategoryRelationsIDs: %w", err)
// 		}

// 		return nil
// 	})

// 	if txErr != nil {
// 		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(txErr))
// 		return
// 	}

// 	illustration := service.FetchRelationInfoForIllustrations(c, s.Store, image)

// 	c.JSON(http.StatusOK, gin.H{
// 		"illustration": illustration,
// 		"message":      "illustrationの編集に成功しました",
// 	})
// }

// // DeleteIllustration godoc
// // @Summary Delete an illustration
// // @Description Deletes a specific illustration by its ID.
// // @Tags illustrations
// // @Accept  json
// // @Produce  json
// // @Param   id   path   int  true  "ID of the illustration to delete"
// // @Success 200 {object} gin/H "Returns a success message indicating the illustration has been deleted"
// // @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: Error parsing the 'id' from path parameters"
// // @Failure 404 {object} request/JSONResponse{data=string} "Not Found: No illustration found with the given ID"
// // @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: Failed to delete the illustration due to a server error"
// // @Router /api/v1/admin/illustrations/{id} [delete]
// func (s *Server) DeleteIllustration(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
// 		return
// 	}

// 	// TODO: illustrationとして取得すれば、冗長な関数を削除できそう
// 	image, err := s.Store.GetImage(c, int64(id))
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
// 			return
// 		}

// 		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to GetImage : %w", err)))
// 		return
// 	}

// 	txErr := s.Store.ExecTx(c.Request.Context(), func(q *db.Queries) error {
// 		err = service.DeleteImageSrc(c, &s.Config, image.OriginalSrc)
// 		if err != nil {
// 			return fmt.Errorf("failed to DeleteImageSrc: %w", err)
// 		}

// 		err = service.DeleteImageSrc(c, &s.Config, image.SimpleSrc.String)
// 		if err != nil {
// 			return fmt.Errorf("failed to DeleteImageSrc: %w", err)
// 		}

// 		// TODO: illustrationとして取得できれば、このrelation取得の処理削除できる
// 		// image_child_category_relationsを削除
// 		err = s.Store.DeleteAllImageChildCategoryRelationsByImageID(c, image.ID)
// 		if err != nil {
// 			return fmt.Errorf("failed to DeleteAllImageChildCategoryRelationsByImageID: %w", err)
// 		}

// 		// image_parent_category_relationsを削除
// 		err = s.Store.DeleteAllImageParentCategoryRelationsByImageID(c, image.ID)
// 		if err != nil {
// 			return fmt.Errorf("failed to DeleteAllImageParentCategoryRelationsByImageID: %w", err)
// 		}

// 		// image_character_relationsを削除
// 		err = s.Store.DeleteAllImageCharacterRelationsByImageID(c, image.ID)
// 		if err != nil {
// 			return fmt.Errorf("failed to DeleteAllImageCharacterRelationsByImageID: %w", err)
// 		}

// 		// Imageを削除
// 		err = s.Store.DeleteImage(c, image.ID)
// 		if err != nil {
// 			return fmt.Errorf("failed to DeleteImage: %w", err)
// 		}

// 		return nil
// 	})

// 	if txErr != nil {
// 		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(txErr))
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "illustrationの削除に成功しました",
// 	})
// }
