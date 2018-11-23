package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/thinmonkey/user-manager/router/middlewares"
	"github.com/thinmonkey/user-manager/handle"
)

// Load loads the middlewares, routes, handlers.
func Load(router *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddlerware())
	router.Use(mw...)
	// 404 Handler.
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	router.GET("/healthcheck", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"code": 200})
	})

	digitalKeyGroup := router.Group("/api/v1/")
	{
		//获取车的某个数字钥匙详情
		digitalKeyGroup.GET("user/:userId", handle.GetUser)
		//删除车的某个数字钥匙
		digitalKeyGroup.DELETE("user/:userId", handle.DeleteUser)
		//更新车的某个数字钥匙
		digitalKeyGroup.PUT("user/:userId", handle.UpdateUser)
		//创建车的某个数字钥匙
		digitalKeyGroup.POST("user", handle.CreateUser)
	}

	return router
}
