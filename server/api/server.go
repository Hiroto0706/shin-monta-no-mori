package api

import (
	"fmt"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/pkg/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	store  *db.Store
	router *gin.Engine
	// tokenMaker token.Maker
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store, config util.Config) (*Server, error) {
	// token, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot create token maker : %w", err)
	// }
	server := &Server{
		config: config,
		store:  store,
		// tokenMaker: token,
	}

	router := gin.Default()
	router.Use(CORSMiddleware())

	// // MasterUserの作成
	// err = insertMasterUser(server)
	// if err != nil {
	// 	return nil, err
	// }

	// Userサイドのルート設定
	SetRouters(router)
	server.router = router

	return server, nil
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
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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
	if server.router == nil {
		return fmt.Errorf("router is nil")
	}
	return server.router.Run(address)
}
