package routes

import "github.com/gin-gonic/gin"

func AuthRoutes(r *gin.Engine) {
	r.POST("/signup", func(ctx *gin.Context) {})
	r.POST("/login", func(ctx *gin.Context) {})
}
