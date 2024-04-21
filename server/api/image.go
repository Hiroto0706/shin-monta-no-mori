package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/internal/domains/service"
	"shin-monta-no-mori/server/pkg/util"
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
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(fmt.Errorf("failed to parse page number from query param : %w", err)))
	}

	illustrations := []*model.Illustration{}

	arg := db.ListImageParams{
		Limit:  int32(server.Config.ImageFetchLimit),
		Offset: int32(page * server.Config.ImageFetchLimit),
	}
	images, err := server.Store.ListImage(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse(fmt.Errorf("failed to ListImage() : %w", err)))
		}

		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to ListImage() : %w", err)))
	}

	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(c, server.Store, i)

		illustrations = append(illustrations, il)
	}

	c.JSON(http.StatusOK, illustrations)
}
