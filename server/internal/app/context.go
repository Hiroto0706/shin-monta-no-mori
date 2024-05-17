package app

import (
	"net/http"

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
func (ctx *AppContext) ErrorResponse(err error) (int, gin.H) {
	return http.StatusInternalServerError, gin.H{"error": err.Error()}
}
