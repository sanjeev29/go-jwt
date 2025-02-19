package routes

import "github.com/gin-gonic/gin"

func UserRoutes(r *gin.Engine) {
	// TODO: Use auth middleware

	r.GET("/users", func(ctx *gin.Context) {})
	r.GET("/users/:user_id", func(ctx *gin.Context) {})
}
