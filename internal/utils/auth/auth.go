package auth

import (
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/gin-gonic/gin"
)

func GetCurrentUser(ctx *gin.Context) (model.User, bool) {
	u, exists := ctx.Get("user")
	if !exists {
		return model.User{}, false
	}
	user, ok := u.(model.User)
	return user, ok
}
