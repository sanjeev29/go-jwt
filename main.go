package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sanjeev29/go-jwt/database"
	"github.com/sanjeev29/go-jwt/routes"
)

func main() {

	// Logger and Recovery Middleware attached.
	router := gin.Default()

	// Add custom routes
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// Connect to database
	_ = database.DBInstance()

	// Router handler functions
	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"success": "access granted for api-1",
		})
	})

	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"success": "access granted for api-2",
		})
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run()
}
