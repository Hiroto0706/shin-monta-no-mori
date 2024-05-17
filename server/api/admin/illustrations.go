package admin

import (
	"fmt"
	"net/http"
	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/internal/domains/service"
	"shin-monta-no-mori/server/pkg/util"
)

const (
	IMAGE_TYPE_IMAGE = "image"
)

type listIllustrationsRequest struct {
	Page int64 `form:"p"`
}

// ListIllustrations godoc
// @Summary List illustrations
// @Description Retrieves a paginated list of illustrations based on the provided page number.
// @Accept  json
// @Produce  json
// @Param   p   query   int  true  "Page number for pagination"
// @Success 200 {array} model/Illustration "A list of illustrations"
// @Failure 400 {object} request/JSONResponse{data=string} "Bad Request: The request is malformed or missing required fields."
// @Failure 500 {object} request/JSONResponse{data=string} "Internal Server Error: An error occurred on the server which prevented the completion of the request."
// @Router /api/v1/admin/illustrations/list [get]
func ListIllustrations(c *app.AppContext) {
	// TODO: bind 周りの処理は関数化して共通化したほうがいい
	var req listIllustrationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(err))
		return
	}

	illustrations := []*model.Illustration{}

	arg := db.ListImageParams{
		Limit:  int32(c.Server.Config.ImageFetchLimit),
		Offset: int32(int(req.Page) * c.Server.Config.ImageFetchLimit),
	}
	images, err := c.Server.Store.ListImage(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse(fmt.Errorf("failed to ListImage : %w", err)))
		return
	}

	for _, i := range images {
		il := service.FetchRelationInfoForIllustrations(c.Context, c.Server.Store, i)

		illustrations = append(illustrations, il)
	}

	c.JSON(http.StatusOK, illustrations)
}
