package routes

import (
	"ScArium/common/config"
	"ScArium/common/log"
	"ScArium/internal/backend/database/entity"
	"ScArium/internal/backend/database/repo"
	"ScArium/internal/backend/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func InitAuthRoutes(r *gin.RouterGroup) {
	r.GET("login", loginRoute())
	r.POST("register", registerRoute(false))
	r.POST("adminRegister", registerRoute(true))
}

func loginRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad json data"})
			c.Abort()
			return
		}
		user, err := repo.GetUserByName(request.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User or password is invalid"})
			c.Abort()
			return
		}
		if !helper.CheckPasswordHash(request.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User or password is invalid"})
			c.Abort()
			return
		}
		jwt, err := helper.GenerateJWT(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate a jwt token"})
			c.Abort()
			return
		}
		c.SetCookie("jwt", jwt, 60*60*24*30, "/", "localhost", false, true)
		c.JSON(200, gin.H{"message": "successfully logged in"})
	}
}

func registerRoute(isAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad json data"})
			c.Abort()
			return
		}
		if repo.DoesUserByNameExists(request.Username) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
			c.Abort()
			return
		}
		hashedPassword, err := helper.HashPassword(request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
			c.Abort()
			return
		}
		user := entity.User{
			Username: request.Username,
			Password: hashedPassword,
			Admin:    isAdmin,
		}
		err = repo.InsertNewUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new user"})
			c.Abort()
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Account successfully created"})
		if isAdmin {
			_, err := os.Create(config.ADMIN_REGISTER_EXIST)
			if err != nil {
				log.I.Fatal("Could not create admin register exist file")
				return
			}
			config.RT_ADMIN_REGISTER_AVAILABLE = false
		}
	}
}
