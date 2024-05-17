package api

// import (
// 	"fmt"
// 	db "shin-monta-no-mori/server/internal/db/sqlc"
// 	"shin-monta-no-mori/server/pkg/token"
// 	"shin-monta-no-mori/server/pkg/util"

// 	"github.com/gin-gonic/gin"
// )

// // Server は、アプリケーション全体の設定、依存関係、およびルーターを保持する構造体
// type Server struct {
// 	Config     util.Config
// 	Store      *db.Store
// 	Router     *gin.Engine
// 	TokenMaker token.Maker
// }

// // NewServer creates a new HTTP server and setup routing
// func NewServer(store *db.Store, config util.Config) (*Server, error) {
// 	token, err := token.NewPasetoMaker(config.TokenSymmetricKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot create token maker : %w", err)
// 	}
// 	server := &Server{
// 		Config:     config,
// 		Store:      store,
// 		TokenMaker: token,
// 	}

// 	router := gin.Default()
// 	router.Use(CORSMiddleware())
// 	server.Router = router

// 	// Userサイドのルート設定
// 	SetUserRouters(server)
// 	// Adminサイドのルート設定
// 	SetAdminRouters(server)

// 	return server, nil
// }

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// TODO: local, stg, prdで値を変更する
// 		allowedOrigins := []string{"http://localhost:3000", "http://localhost:3030"}
// 		origin := c.GetHeader("Origin")

// 		for _, allowedOrigin := range allowedOrigins {
// 			if origin == allowedOrigin {
// 				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 				c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 				c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, multipart/form-data")
// 				c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
// 				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

// 				break
// 			}
// 		}

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

// func (server *Server) Start(address string) error {
// 	if server.Router == nil {
// 		return fmt.Errorf("router is nil")
// 	}
// 	return server.Router.Run(address)
// }
