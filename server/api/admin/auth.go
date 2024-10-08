package admin

import (
	"database/sql"
	"net/http"

	"shin-monta-no-mori/internal/app"
	db "shin-monta-no-mori/internal/db/sqlc"
	"shin-monta-no-mori/pkg/lib/password"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		ctx.Server.Logger.Error("failed to GetOperatorByEmail",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(err))
		return
	}

	if err = password.CheckPassword(req.Password, operator.HashedPassword); err != nil {
		ctx.Server.Logger.Info("failed to CheckPassword",
			zap.Error(err),
		)
		ctx.JSON(http.StatusUnauthorized, app.ErrorResponse(err))
		return
	}
	if err = password.CheckEmail(req.Email, operator.Email); err != nil {
		ctx.Server.Logger.Info("failed to CheckEmail",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		ctx.JSON(http.StatusUnauthorized, app.ErrorResponse(err))
		return
	}

	// アクセストークンとリフレッシュトークンを作成し、セッションに保存
	accessToken, _, err := ctx.Server.TokenMaker.CreateToken(
		operator.Name,
		ctx.Server.Config.AccessTokenDuration,
	)
	if err != nil {
		ctx.Server.Logger.Info("failed to CreateToken",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := ctx.Server.TokenMaker.CreateToken(
		operator.Name,
		ctx.Server.Config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.Server.Logger.Info("failed to CreateToken",
			zap.String("email", req.Email),
			zap.Error(err),
		)
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
		ctx.Server.Logger.Info("failed to CreateSessionParams",
			zap.String("operator_name", operator.Name),
			zap.String("email", req.Email),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, app.ErrorResponse(err))
		return
	}

	rsp := loginResponse{
		AccessToken: accessToken,
	}

	ctx.Server.Logger.Info("login success",
		zap.Int64("operator_id", operator.ID),
		zap.String("operator_name", operator.Name),
		zap.String("email", req.Email),
		zap.Error(err),
	)

	ctx.JSON(http.StatusOK, rsp)
}

type verifyRequest struct {
	AccessToken string `form:"access_token" binding:"required"`
}

func VerifyAccessToken(ctx *app.AppContext) {
	var req verifyRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Server.Logger.Info("failed to bind accessToken", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}
	_, err := ctx.Server.TokenMaker.VerifyToken(req.AccessToken)
	if err != nil {
		ctx.Server.Logger.Info("failed to VerifyToken", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, app.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": true})
}
