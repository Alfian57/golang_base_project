package middleware

import (
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/logger"
	"github.com/Alfian57/belajar-golang/internal/response"
	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Errorw("panic recovered", "error", err)
				response.WriteErrorResponse(c, errs.ErrInternalServer)
				c.Abort()
			}
		}()
		c.Next()
	}
}
