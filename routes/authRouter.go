package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sanjeev29/go-jwt/controllers"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/signup", controller.Signup())
	r.POST("/login", func(ctx *gin.Context) {})
}
