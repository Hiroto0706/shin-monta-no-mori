package util

import "github.com/gin-gonic/gin"

func C(f func()) gin.HandlerFunc {
	return func(c *gin.Context) {
		return
	}
}
