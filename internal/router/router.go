package router

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	api := router.Group("api")

	v1 := api.Group("v1")
	RegisterV1Route(v1)

	return router
}
