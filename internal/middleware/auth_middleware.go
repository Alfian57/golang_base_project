package middleware

import (
	"github.com/Alfian57/belajar-golang/internal/di"
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/response"
	"github.com/Alfian57/belajar-golang/internal/utils/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authService := di.InitializeUserService()

		accessToken, err := ctx.Cookie("access_token")
		if err != nil {
			response.WriteErrorResponse(ctx, errs.ErrUnauthorized)
			ctx.Abort()
			return
		}

		userID, err := jwt.ValidateAccessToken(accessToken)
		if err != nil {
			response.WriteErrorResponse(ctx, errs.ErrUnauthorized)
			ctx.Abort()
			return
		}

		if userID == "" {
			response.WriteErrorResponse(ctx, errs.ErrUnauthorized)
			ctx.Abort()
			return
		}

		user, err := authService.GetUserByID(ctx, userID)
		if err != nil {
			response.WriteErrorResponse(ctx, errs.ErrUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set("access_token", accessToken)
		ctx.Set("user", user)

		ctx.Next()
	}
}
