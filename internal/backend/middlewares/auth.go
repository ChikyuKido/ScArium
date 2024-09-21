package middlewares

import (
	"ScArium/internal/backend/database/repo"
	"ScArium/internal/backend/helper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, _ := c.Cookie("jwt")
		// guest login
		if tokenString == "" || tokenString == "guest" {
			guest, err := repo.GetUserByName("guest")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No guest user found."})
				c.Abort()
				return
			}
			c.Set("user", guest)
			c.Next()
			return
		}

		token, err := helper.GetToken(tokenString)
		if err != nil || !token.Valid {
			c.SetCookie("jwt", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.SetCookie("jwt", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		userName := claims["username"].(string)
		user, err := repo.GetUserByName(userName)
		if err != nil {
			c.SetCookie("jwt", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
