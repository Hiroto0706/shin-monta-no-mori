package api

import (
	admin_handlers "shin-monta-no-mori/server/api/handlers/admin"
	user_handlers "shin-monta-no-mori/server/api/handlers/user"

	"github.com/gin-gonic/gin"
)

func SetUserRouters(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", user_handlers.Greet)
	}
}

func SetAdminRouters(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	admin := v1.Group("/admin")
	// ログイン認証必須
	// admin.Use(authMiddleware(server.tokenMaker))
	{
		admin.GET("/", admin_handlers.Greet)
	}
}
