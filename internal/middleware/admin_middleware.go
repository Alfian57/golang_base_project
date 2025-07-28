package middleware

import (
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/Alfian57/belajar-golang/internal/response"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(model.User)

		if user.Role != model.UserRoleAdmin {
			response.WriteErrorResponse(ctx, errs.ErrForbidden)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
