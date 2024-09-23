package middlewares

import (
	"ScArium/common/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AdminRegisterRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.RT_ADMIN_REGISTER_AVAILABLE {
			if c.FullPath() == "/api/v1/auth/adminRegister" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "The admin register is disabled"})
				c.Abort()
				return
			}
			if c.FullPath() == "/auth/adminRegister" {
				c.Redirect(http.StatusTemporaryRedirect, "/auth/register")
				c.Abort()
				return
			}
			return
		}
		if c.FullPath() == "/auth/adminRegister" {
			return
		}
		if c.FullPath() == "/api/v1/auth/adminRegister" {
			return
		}
		if strings.HasPrefix(c.FullPath(), "/api") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The endpoints are disabled until a admin user was registered"})
			c.Abort()
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, "/auth/adminRegister")
		c.Abort()
	}
}
