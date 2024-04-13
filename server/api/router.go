package api

import (
	"shin-monta-no-mori/server/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetRouters(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", handlers.Greet)
	}
}
