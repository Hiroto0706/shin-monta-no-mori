package app

import (
	"github.com/gin-gonic/gin"
)

// AppContext は gin.Context を拡張したコンテキスト
type AppContext struct {
	*gin.Context
	Server *Server
}

// NewAppContext は新しい AppContext を作成
func NewAppContext(c *gin.Context, s *Server) *AppContext {
	return &AppContext{
		Context: c,
		Server:  s,
	}
}

// ErrorResponse はエラーレスポンスを生成
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
