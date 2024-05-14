package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"shin-monta-no-mori/server/pkg/token"
	"shin-monta-no-mori/server/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 認証ヘッダーを取得
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// ヘッダーが存在しない場合、エラーを返す
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.NewErrorResponse(err))
			return
		}

		// ヘッダーをスペースで分割して、認証タイプとトークンを取得
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.NewErrorResponse(err))
			return
		}

		// 認証タイプを確認（ここでは "bearer" のみをサポート）
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.NewErrorResponse(err))
			return
		}

		// トークンを検証
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.NewErrorResponse(err))
			return
		}

		// トークンのペイロードをコンテキストに保存して、次のハンドラに進む
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
