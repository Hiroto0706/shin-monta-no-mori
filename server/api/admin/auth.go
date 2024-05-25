package admin

import (
	"database/sql"
	"net/http"
	"time"

	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/pkg/util"

	"github.com/google/uuid"
)

type (
	operatorResponse struct {
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}

	loginRequest struct {
		Username string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}

	loginResponse struct {
		SessionID             uuid.UUID        `json:"session_id"`
		AccessToken           string           `json:"access_token"`
		AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
		RefreshToken          string           `json:"refresh_token"`
		RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
		Operator              operatorResponse `json:"operator"`
	}
)

func Login(ctx *app.AppContext) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ErrorResponse(err))
		return
	}

	operator, err := ctx.Server.Store.GetOperatorByName(ctx, req.Username)
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
	accessToken, accessPayload, err := ctx.Server.TokenMaker.CreateToken(
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

	session, err := ctx.Server.Store.CreateSession(ctx, db.CreateSessionParams{
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
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		Operator: operatorResponse{
			Email:     operator.Email,
			Username:  operator.Name,
			CreatedAt: operator.CreatedAt,
		},
	}

	ctx.JSON(http.StatusOK, rsp)
}
