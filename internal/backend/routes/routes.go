package routes

import (
	"ScArium/internal/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	adminRegisterRedirectGroup := r.Group("/", middlewares.AdminRegisterRedirectMiddleware())
	backend := adminRegisterRedirectGroup.Group("/api/v1")
	backendAuth := adminRegisterRedirectGroup.Group("/api/v1", middlewares.AuthMiddleware(false))
	//to check if the user is logged in
	backendAuth.GET("/auth/pingAuth", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	InitAuthRoutes(backend.Group("/auth"))
	InitAccountRoutes(backendAuth.Group("/account"))
	InitSitesRoutes(r)
	initDefaultRoutes(r)
}

func initDefaultRoutes(r *gin.Engine) {
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"error": "Method Not Allowed",
		})
	})
}
