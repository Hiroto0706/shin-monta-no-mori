package middleware_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shin-monta-no-mori/server/api/middleware"
	"shin-monta-no-mori/server/pkg/token"
	"shin-monta-no-mori/server/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func TestAuthMiddleware(t *testing.T) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	require.NoError(t, err)

	// テストケースを定義
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "正常系",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				accessToken, _, err := tokenMaker.CreateToken("testuser", time.Minute)
				require.NoError(t, err)
				request.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessToken)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "異常系（NoAuthorization）",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// 認証ヘッダーを設定しない
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "異常系（InvalidAuthorizationFormat）",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				request.Header.Set(authorizationHeaderKey, "invalid_format")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "異常系（UnsupportedAuthorizationType）",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				accessToken, _, err := tokenMaker.CreateToken("testuser", time.Minute)
				require.NoError(t, err)
				request.Header.Set(authorizationHeaderKey, "basic "+accessToken)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "異常系（ExpiredToken）",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				accessToken, _, err := tokenMaker.CreateToken("testuser", -time.Minute)
				require.NoError(t, err)
				request.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessToken)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.AuthMiddleware(tokenMaker))
			router.GET("/api/v1/admin/illustrations/list", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"status": "success"})
			})

			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/illustrations/list", nil)

			// 認証ヘッダーを設定
			if tc.setupAuth != nil {
				tc.setupAuth(t, request, tokenMaker)
			}

			// リクエストを送信
			router.ServeHTTP(recorder, request)

			// レスポンスを確認
			tc.checkResponse(t, recorder)
		})
	}
}
