package routes

import (
	"ScArium/common/log"
	"ScArium/internal/backend/database/repo"
	"ScArium/internal/backend/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitAccountRoutes(r *gin.RouterGroup) {
	r.POST("createD4sAccount", createD4sAccount())
	r.POST("createMoodleAccount", createMoodleAccount())
	r.GET("getMoodleAccounts", getMoodleAccounts())
	r.GET("getD4sAccounts", getD4sAccounts())
}

func getMoodleAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helper.GetUserFromContext(c)
		if user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The authentication failed"})
			log.I.Fatal("The user does not exist in a authenticated endpoint. This endpoint should never have been reached")
			return
		}
		accounts, err := repo.GetMoodleAccounts(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve moodle accounts for user"})
			return
		}
		c.JSON(http.StatusOK, accounts)
	}
}

func getD4sAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helper.GetUserFromContext(c)
		if user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The authentication failed"})
			log.I.Fatal("The user does not exist in a authenticated endpoint. This endpoint should never have been reached")
			return
		}
		accounts, err := repo.GetD4sAccounts(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve d4s accounts for user"})
			return
		}
		c.JSON(http.StatusOK, accounts)
	}
}

func createMoodleAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			InstanceUrl string `json:"instance_url" binding:"required"`
			Username    string `json:"username" binding:"required"`
			Password    string `json:"password" binding:"required"`
			DisplayName string `json:"display_name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad json format"})
			return
		}
		user := helper.GetUserFromContext(c)
		if user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The authentication failed"})
			log.I.Fatal("The user does not exist in a authenticated endpoint. This endpoint should never have been reached")
			return
		}
		err := repo.CreateMoodleAccount(user, request.InstanceUrl, request.Username, request.Password, request.DisplayName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create the account"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Moodle Account created successfully"})
	}
}
func createD4sAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Username    string `json:"username" binding:"required"`
			Password    string `json:"password" binding:"required"`
			DisplayName string `json:"display_name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad json format"})
			return
		}
		user := helper.GetUserFromContext(c)
		if user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The authentication failed"})
			log.I.Fatal("The user does not exist in a authenticated endpoint. This endpoint should never have been reached")
			return
		}
		err := repo.CreateD4sAccount(user, request.Username, request.Password, request.DisplayName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create the account"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Digi4School Account created successfully"})
	}
}
