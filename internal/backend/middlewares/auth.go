package middlewares

import (
	"ScArium/internal/backend/database/repo"
	"ScArium/internal/backend/helper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, _ := c.Cookie("jwt")
		if tokenString == "" {
			if redirect {
				c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No jwt token provided"})
				c.Abort()
			}
			return
		}

		token, err := helper.GetToken(tokenString)
		if err != nil || !token.Valid {
			c.SetCookie("jwt", "", -1, "/", "", false, true)
			if redirect {
				c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
			}
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.SetCookie("jwt", "", -1, "/", "", false, true)
			if redirect {
				c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
			}
			return
		}
		userName := claims["username"].(string)
		user, err := repo.GetUserByName(userName)
		if err != nil {
			c.SetCookie("jwt", "", -1, "/", "", false, true)
			if redirect {
				c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
			}
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
