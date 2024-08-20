package app

import (
	"fmt"
	"net/http"

	db "shin-monta-no-mori/internal/db/sqlc"
	"shin-monta-no-mori/pkg/token"
	"shin-monta-no-mori/pkg/util"

	"github.com/gin-gonic/gin"
)

// Server は、アプリケーション全体の設定、依存関係、およびルーターを保持する構造体
type Server struct {
	Config      util.Config
	Store       *db.Store
	RedisClient RedisClient
	Router      *gin.Engine
	TokenMaker  token.Maker
}

// NewServer は新しいサーバーインスタンスを作成
func NewServer(config util.Config, store *db.Store, redis RedisClient, tokenMaker token.Maker) *Server {
	server := &Server{
		Config:      config,
		Store:       store,
		RedisClient: redis,
		Router:      gin.Default(),
		TokenMaker:  tokenMaker,
	}

	router := gin.Default()
	server.Router = router

	return server
}

func CORSMiddleware(config util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := config.Origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (server *Server) Start(address string) error {
	if server.Router == nil {
		return fmt.Errorf("router is nil")
	}
	return server.Router.Run(address)
}
