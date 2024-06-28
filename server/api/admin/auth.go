package admin

import (
	"database/sql"
	"net/http"

	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/pkg/util"

	"github.com/gin-gonic/gin"
)

type (
	loginRequest struct {
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}

	loginResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func Login(ctx *app.AppContext) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	operator, err := ctx.Server.Store.GetOperatorByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, app.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(err))
		return
	}

	if err = util.CheckPassword(req.Password, operator.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, app.ErrorResponse(err))
		return
	}
	if err = util.CheckEmail(req.Email, operator.Email); err != nil {
		ctx.JSON(http.StatusUnauthorized, app.ErrorResponse(err))
		return
	}

	// アクセストークンとリフレッシュトークンを作成し、セッションに保存
	accessToken, _, err := ctx.Server.TokenMaker.CreateToken(
		operator.Name,
		ctx.Server.Config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := ctx.Server.TokenMaker.CreateToken(
		operator.Name,
		ctx.Server.Config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(err))
		return
	}

	_, err = ctx.Server.Store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Name:         operator.Name,
		Email:        sql.NullString{String: req.Email, Valid: true},
		RefreshToken: refreshToken,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(err))
		return
	}

	rsp := loginResponse{
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type verifyRequest struct {
	AccessToken string `form:"access_token" binding:"required"`
}

func VerifyAccessToken(ctx *app.AppContext) {
	var req verifyRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	_, err := ctx.Server.TokenMaker.VerifyToken(req.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, app.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": true})
}
