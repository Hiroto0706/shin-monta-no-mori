package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) Greet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "new Hello World from server.",
	})
}

func (server *Server) ListIllustrations(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("p"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Errorf("failed to parse page number from query param : %w", err)))
		return
	}
	log.Println(page)

	illustrations := []model.Illustration{}
	log.Println(server.Config.ImageFetchLimit)

	arg := db.ListImageParams{
		Limit:  int32(server.Config.ImageFetchLimit),
		Offset: int32(page * server.Config.ImageFetchLimit),
	}
	log.Println(arg)
	images, err := server.Store.ListImage(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, NewErrorResponse(fmt.Errorf("failed to ListImage() : %w", err)))
			return
		}

		c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("failed to ListImage() : %w", err)))
		return
	}
	log.Println(images)
	for _, i := range images {
		il := model.Illustration{}
		il.Image = i
		illustrations = append(illustrations, il)
	}

	log.Println(illustrations)

	// resImages := []responseImage{}
	// for _, image := range images {
	// 	resImage := responseImage{}
	// 	resImage.Image = image

	// 	typeID := resImage.Image.TypeID
	// 	imageType, err := server.store.GetType(c, typeID)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("failed to get type id from image type_id : %w", err)))
	// 		return
	// 	}
	// 	resImage.ImageType = imageType

	// 	imageID := resImage.Image.ID
	// 	imageCategories, err := server.store.ListImageCategoriesByImage(c, imageID)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("failed to get image_categories by image id : %w", err)))
	// 		return
	// 	}
	// 	categories := []db.Category{}
	// 	for _, imageCategory := range imageCategories {
	// 		category, err := server.store.GetCategory(c, imageCategory.CategoryID)
	// 		if err != nil {
	// 			c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("failed to get category : %w", err)))
	// 			return
	// 		}
	// 		categories = append(categories, category)
	// 	}
	// 	resImage.Categories = categories
	// 	resImages = append(resImages, resImage)
	// }

	c.JSON(http.StatusOK, gin.H{
		"illustrations": illustrations,
	})
}
