package routes

import (
	"ScArium/internal/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	adminRegisterRedirectGroup := r.Group("/", middlewares.AdminRegisterRedirectMiddleware())
	backend := adminRegisterRedirectGroup.Group("/api/v1")
	InitAuthRoutes(backend.Group("/auth"))
	initDefaultRoutes(r)
}

func initDefaultRoutes(r *gin.Engine) {
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"error": "Method Not Allowed",
		})
	})
}
