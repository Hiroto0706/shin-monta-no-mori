package app

import (
	"fmt"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/pkg/token"
	"shin-monta-no-mori/server/pkg/util"

	"github.com/gin-gonic/gin"
)

// Server は、アプリケーション全体の設定、依存関係、およびルーターを保持する構造体
type Server struct {
	Config     util.Config
	Store      *db.Store
	Router     *gin.Engine
	TokenMaker token.Maker
}

// NewServer は新しいサーバーインスタンスを作成します。
func NewServer(config util.Config, store *db.Store, tokenMaker token.Maker) *Server {
	server := &Server{
		Config:     config,
		Store:      store,
		Router:     gin.Default(),
		TokenMaker: tokenMaker,
	}

	router := gin.Default()
	router.Use(CORSMiddleware())
	server.Router = router

	return server
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: local, stg, prdで値を変更する
		allowedOrigins := []string{"http://localhost:3000", "http://localhost:3030"}
		origin := c.GetHeader("Origin")

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, multipart/form-data")
				c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

				break
			}
		}

		if c.Request.Method == "OPTIONS" {
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
