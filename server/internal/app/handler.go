package app

import "github.com/gin-gonic/gin"

// HandlerFuncWrapper は、AppContextを使用する汎用的な中継関数です。
func HandlerFuncWrapper(s *Server, handler func(*AppContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewAppContext(c, s)
		handler(ctx)
	}
}
