package binder

import (
	"fmt"
	"net/http"
	"shin-monta-no-mori/server/internal/app"

	"github.com/gin-gonic/gin"
)

// BindQuery クエリパラメータと構造体をバインドする
func BindQuery(ctx *gin.Context, req interface{}) error {
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(fmt.Errorf("failed to c.ShouldBindQuery : %w", err)))
		return err
	}
	return nil
}
