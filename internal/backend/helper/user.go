package helper

import (
	"ScArium/internal/backend/database/entity"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(ctx *gin.Context) *entity.User {
	value, exists := ctx.Get("user")
	if !exists {
		return nil
	}
	user, ok := value.(*entity.User)
	if !ok {
		return nil
	}
	return user
}
